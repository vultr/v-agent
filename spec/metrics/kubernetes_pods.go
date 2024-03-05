// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"
	"github.com/vultr/v-agent/spec/connectors"
	"github.com/vultr/v-agent/spec/wrkld"

	"go.uber.org/zap"
)

// ProbeKubernetesPod pulls metrics from a specific pod
func ProbeKubernetesPod(endpoint, port, path string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://%s:%s%s", endpoint, port, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ScrapeKubernetesPods scrapes /metrics of all pods in specified namespaces that have metric collection enabled
func ScrapeKubernetesPods() error {
	log := zap.L().Sugar()

	clientset, err := connectors.GetKubernetesConn()
	if err != nil {
		return err
	}

	namespaces := config.GetKubernetesPodsNamespaces()
	for i := range namespaces {
		pods, err := wrkld.GetScrapeablePods(clientset, namespaces[i])
		if err != nil {
			log.Error(err)

			continue
		}

		for j := range pods {
			annotations := pods[j].GetAnnotations()

			annoPort, ok := annotations["prometheus.io/port"]
			if !ok {
				log.With(
					"pod", pods[j].ObjectMeta.Name,
				).Error("prometheus.io/port does not exist on pod")
			}

			annoPath, ok := annotations["prometheus.io/path"]
			if !ok {
				log.With(
					"pod", pods[j].ObjectMeta.Name,
				).Warn("prometheus.io/path does not exist on pod using /metrics")

				annoPath = "/metrics"
			}

			podIP := pods[j].Status.PodIP

			log.Infof("scraping pod %s (namespace=%s)", pods[j].ObjectMeta.Name, pods[j].ObjectMeta.Namespace)

			data, err := ProbeKubernetesPod(podIP, annoPort, annoPath)
			if err != nil {
				log.With(
					"pod", pods[j].ObjectMeta.Name,
				).Warnf("error scraping pod metrics: %s", err.Error())

				continue
			}

			podMetrics, err := parseMetrics(data)
			if err != nil {
				log.With(
					"pod", pods[j].ObjectMeta.Name,
				).Warnf("error parsing pod metrics: %s", err.Error())

				continue
			}

			tsList := GetMetricsAsTimeSeries(podMetrics)

			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}

			var ba *BasicAuth
			if cfg.BasicAuthUser != "" && cfg.BasicAuthPass != "" {
				ba = &BasicAuth{
					Username: cfg.BasicAuthUser,
					Password: cfg.BasicAuthPass,
				}
			}

			wc, err := NewWriteClient(cfg.Endpoint, &HTTPConfig{
				Timeout:   5 * time.Second,
				BasicAuth: ba,
			})
			if err != nil {
				return err
			}

			log.Infof("sending metrics for pod %s", pods[j].ObjectMeta.Name)

			if err := wc.Store(context.Background(), tsList); err != nil {
				return err
			}
		}
	}

	return nil
}
