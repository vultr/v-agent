// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"github.com/levigross/grequests"
)

var (
	// ErrEtcdUnhealthy returned if response is not "true" from /health
	ErrEtcdUnhealthy = errors.New("kube-apiserver unhealthy")
)

// HealthResp response for /health is marshalled to this
type HealthResp struct {
	Health string `json:"health"`
	Reason string `json:"reason"`
}

// DoEtcdHealthCheck probes /health and returns nil or ErrKubeApiServerUnhealthy, or some other error
func DoEtcdHealthCheck() error {
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

	caCertData, _ := ioutil.ReadFile(*caCert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	certPair, _ := tls.LoadX509KeyPair(*cert, *key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certPair},
			},
		},
		Timeout: 5 * time.Second,
	}

	endpoint, err := config.GetEtcdEndpoint()
	if err != nil {
		return err
	}

	resp, err := grequests.Get(fmt.Sprintf("%s/health", *endpoint),
		&grequests.RequestOptions{
			HTTPClient: client,
		})
	if err != nil {
		return err
	}

	var jsonResp HealthResp

	if err := resp.JSON(&jsonResp); err != nil {
		return err
	}

	if jsonResp.Health == "true" {
		return nil
	} else {
		return ErrEtcdUnhealthy
	}
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

	caCertData, _ := ioutil.ReadFile(*caCert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	certPair, _ := tls.LoadX509KeyPair(*cert, *key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certPair},
			},
		},
		Timeout: 5 * time.Second,
	}

	endpoint, err := config.GetEtcdEndpoint()
	if err != nil {
		return nil, err
	}

	resp, err := grequests.Get(fmt.Sprintf("%s/metrics", *endpoint),
		&grequests.RequestOptions{
			HTTPClient: client,
		})
	if err != nil {
		return nil, err
	}

	return resp.Bytes(), nil
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
