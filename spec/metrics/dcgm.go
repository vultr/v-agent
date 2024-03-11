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
	v1 "k8s.io/api/core/v1"
)

// ProbeDCGMMetrics probes /metrics from DCGM
func ProbeDCGMMetrics(target string, port int32) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://%s/metrics", target))
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

// ScrapeDCGMMetrics scrapes nvidia DCGM metrics
func ScrapeDCGMMetrics() error {
	log := zap.L().Sugar()

	clientset, err := connectors.GetKubernetesConn()
	if err != nil {
		return err
	}

	ns, err := config.GetDCGMNamespace()
	if err != nil {
		return err
	}

	ep, err := config.GetDCGMEndpoint()
	if err != nil {
		return err
	}

	dcgmEndpoints, err := wrkld.GetEndpoint(clientset, *ns, *ep)
	if err != nil {
		return err
	}

	for i := range dcgmEndpoints.Subsets {
		var port v1.EndpointPort
		var addr v1.EndpointAddress

		for j := range dcgmEndpoints.Subsets[i].Ports {
			if dcgmEndpoints.Subsets[i].Ports[j].Name == "gpu-metrics" { // extract port
				port = dcgmEndpoints.Subsets[i].Ports[j]
			}
		}

		for j := range dcgmEndpoints.Subsets[i].Addresses {
			addr = dcgmEndpoints.Subsets[i].Addresses[j]
		}

		log.With(
			"ip", addr.IP,
			"port", port.Port,
		).Info("attempting to scrape dcgm metrics")

		dcgmResp, err := ProbeDCGMMetrics(addr.IP, port.Port)
		if err != nil {
			log.Error(err)

			continue
		}

		dcgmMetrics, err := parseMetrics(dcgmResp)
		if err != nil {
			log.Error(err)

			continue
		}

		tsList := GetMetricsAsTimeSeries(dcgmMetrics)

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

		if err := wc.Store(context.Background(), tsList); err != nil {
			return err
		}
	}

	return nil
}
