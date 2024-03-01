// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"
)

// HealthResp response for /health is marshaled to this
type HealthResp struct {
	Health string `json:"health"`
	Reason string `json:"reason"`
}

// DoEtcdHealthCheck probes /health and returns nil or ErrKubeApiServerUnhealthy, or some other error
func DoEtcdHealthCheck() error {
	var jsonResp HealthResp

	caCert, err := config.GetEtcdCACert()
	if err != nil {
		return err
	}
	cert, err := config.GetEtcdClientCert()
	if err != nil {
		return err
	}
	key, err := config.GetEtcdClientKey()
	if err != nil {
		return err
	}

	caCertData, _ := os.ReadFile(*caCert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	certPair, _ := tls.LoadX509KeyPair(*cert, *key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ //nolint
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certPair},
			},
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	endpoint, err := config.GetEtcdEndpoint()
	if err != nil {
		return err
	}

	resp, err := client.Get(fmt.Sprintf("%s/health", *endpoint))
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Health == "true" {
		return nil
	}

	return ErrEtcdUnhealthy
}

// ProbeEtcdMetrics probes /metrics from etcd
func ProbeEtcdMetrics() ([]byte, error) {
	caCert, err := config.GetEtcdCACert()
	if err != nil {
		return nil, err
	}
	cert, err := config.GetEtcdClientCert()
	if err != nil {
		return nil, err
	}
	key, err := config.GetEtcdClientKey()
	if err != nil {
		return nil, err
	}

	caCertData, _ := os.ReadFile(*caCert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	certPair, _ := tls.LoadX509KeyPair(*cert, *key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ //nolint
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certPair},
			},
			DisableKeepAlives: true, // very important, prevents connection pooling, which can leak connections
		},
		Timeout: 5 * time.Second,
	}

	endpoint, err := config.GetEtcdEndpoint()
	if err != nil {
		return nil, err
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

// ScrapeEtcdMetrics scrapes kube-apiserver /metrics endpoint and remote writes the metrics
func ScrapeEtcdMetrics() error {
	etcdResp, err := ProbeEtcdMetrics()
	if err != nil {
		return err
	}

	etcdMetrics, err := parseMetrics(etcdResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(etcdMetrics)

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
