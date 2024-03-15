// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"syscall"
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

	resp, err := client.Get(fmt.Sprintf("http://%s:%d/metrics", target, port))
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

	ns := config.GetDCGMNamespace()
	ep := config.GetDCGMEndpoint()

	dcgmEndpoints, err := wrkld.GetEndpoint(clientset, ns, ep)
	if err != nil {
		return err
	}

	cfg := config.GetConfig()

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

	for i := range dcgmEndpoints.Subsets {
		var port v1.EndpointPort

		for j := range dcgmEndpoints.Subsets[i].Ports {
			if dcgmEndpoints.Subsets[i].Ports[j].Name == "gpu-metrics" { // extract port
				port = dcgmEndpoints.Subsets[i].Ports[j]
			}
		}

		// loop each address
		for j := range dcgmEndpoints.Subsets[i].Addresses {
			addr := dcgmEndpoints.Subsets[i].Addresses[j]

			log.Infof("scraping dcgm metrics from %s:%d", addr.IP, port.Port)

			dcgmResp, err := ProbeDCGMMetrics(addr.IP, port.Port)
			if err != nil {
				if errors.Is(err, syscall.ECONNREFUSED) {
					log.Warn(err)
				} else {
					log.Error(err)
				}

				continue
			}

			dcgmMetrics, err := parseMetrics(dcgmResp)
			if err != nil {
				log.Error(err)

				continue
			}

			tsList := GetMetricsAsTimeSeries(dcgmMetrics)

			if err := wc.Store(context.Background(), tsList); err != nil {
				log.Error(err)

				continue
			}
		}
	}

	return nil
}
