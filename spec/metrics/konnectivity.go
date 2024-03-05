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

// DoKonnectivityHealthCheck probes /healthz and returns nil or ErrKonnectivityServerUnhealthy, or some other error
func DoKonnectivityHealthCheck() error {
	endpoint, err := config.GetKonnectivityHealthEndpoint()
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/healthz", *endpoint))
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(data) == "ok" {
		return nil
	}

	return ErrKonnectivityServerUnhealthy
}

// ProbeKonnectivityMetrics probes /metrics from konnectivity
func ProbeKonnectivityMetrics() ([]byte, error) {
	endpoint, err := config.GetKonnectivityMetricsEndpoint()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/metrics", *endpoint))
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

// ScrapeKonnectivityMetrics scrapes konnectivity /metrics endpoint and remote writes the metrics
func ScrapeKonnectivityMetrics() error {
	konnectivityResp, err := ProbeKonnectivityMetrics()
	if err != nil {
		return err
	}

	konnectivityMetrics, err := parseMetrics(konnectivityResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(konnectivityMetrics)

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

	return nil
}
