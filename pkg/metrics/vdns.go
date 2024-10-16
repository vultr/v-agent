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

// DoVDNSHealthCheck probes /metrics and returns nil or ErrVDNSUnhealthy, or some other error
func DoVDNSHealthCheck() error {
	endpoint := config.GetVDNSMetricsEndpoint()

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
		return ErrVDNSUnhealthy
	}

	return nil
}

// ProbeVDNSMetrics probes /metrics from v-dns
func ProbeVDNSMetrics() ([]byte, error) {
	endpoint := config.GetVDNSMetricsEndpoint()

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

// ScrapeVDNSMetrics scrapes v-dns /metrics endpoint and remote writes the metrics
func ScrapeVDNSMetrics() error {
	vdnsResp, err := ProbeVDNSMetrics()
	if err != nil {
		return err
	}

	vdnsMetrics, err := parseMetrics(vdnsResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(vdnsMetrics)

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
