// Package metrics metrics collection
package metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"
)

// DoNginxVTSHealthCheck probes /metrics and returns nil or ErrNginxVTSServerUnhealthy, or some other error
func DoNginxVTSHealthCheck() error {
	endpoint := config.GetNginxVTSMetricsEndpoint()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/metrics", endpoint))
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint

	if _, err := io.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrNginxVTSServerUnhealthy
	}

	return nil
}

// ProbeNginxVTSMetrics probes /metrics from nginx-vts
func ProbeNginxVTSMetrics() ([]byte, error) {
	endpoint := config.GetNginxVTSMetricsEndpoint()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/metrics", endpoint))
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

// ScrapeNginxVTSMetrics scrapes nginx-vts /metrics endpoint and remote writes the metrics
func ScrapeNginxVTSMetrics() error {
	nginxVtsResp, err := ProbeNginxVTSMetrics()
	if err != nil {
		return err
	}

	nginxVtsMetrics, err := parseMetrics(nginxVtsResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(nginxVtsMetrics)

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

	if err := wc.Store(context.Background(), tsList); err != nil {
		return err
	}

	return nil
}
