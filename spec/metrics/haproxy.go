// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"
)

// DoHAProxyHealthCheck probes /metrics and returns nil or ErrHAProxyServerUnhealthy, or some other error
func DoHAProxyHealthCheck() error {
	endpoint := config.GetHAProxyMetricsEndpoint()

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
		return ErrHAProxyServerUnhealthy
	}

	return nil
}

// ProbeHAProxyMetrics probes /metrics from haproxy
func ProbeHAProxyMetrics() ([]byte, error) {
	endpoint := config.GetHAProxyMetricsEndpoint()

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

// ScrapeHAProxyMetrics scrapes haproxy /metrics endpoint and remote writes the metrics
func ScrapeHAProxyMetrics() error {
	haproxyResp, err := ProbeHAProxyMetrics()
	if err != nil {
		return err
	}

	haproxyMetrics, err := parseMetrics(haproxyResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(haproxyMetrics)

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
