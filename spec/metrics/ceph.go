// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"go.uber.org/zap"
)

// ProbeCephMetrics probes /metrics from ceph
func ProbeCephMetrics() ([]byte, error) {
	endpoint, err := config.GetCephMetricsEndpoint()
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

	if len(data) == 0 {
		return nil, ErrCephMgrNotActive
	}

	return data, nil
}

// ScrapeCephMetrics scrapes ceph /metrics endpoint and remote writes the metrics
func ScrapeCephMetrics() error {
	log := zap.L().Sugar()

	cephResp, err := ProbeCephMetrics()
	if err != nil {
		if errors.Is(err, ErrCephMgrNotActive) {
			log.Warn(err)
			return nil
		}

		return err
	}

	// only necessary for broken /metrics implementations
	data := removeMetadata(cephResp)

	cephMetrics, err := parseMetrics(data)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(cephMetrics)

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
