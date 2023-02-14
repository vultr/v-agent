// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// ErrKubeApiServerUnhealthy returned if response is not "ok" from /healthz
	ErrKubeApiServerUnhealthy = errors.New("kube-apiserver unhealthy")
)

// DoKubeApiServerHealthCheck probes /healthz and returns nil or ErrKubeApiServerUnhealthy, or some other error
func DoKubeApiServerHealthCheck() error {
	kubeconfig, err := config.GetKubeconfig()
	if err != nil {
		return err
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	content, err := clientset.Discovery().RESTClient().Get().AbsPath("/healthz").DoRaw(context.TODO())
	if err != nil {
		return err
	}

	if string(content) == "ok" {
		return nil
	} else {
		return ErrKubeApiServerUnhealthy
	}
}

// ProbeKubeApiServerMetrics probes /metrics from kube-apiserver
func ProbeKubeApiServerMetrics() ([]byte, error) {
	kubeconfig, err := config.GetKubeconfig()
	if err != nil {
		return nil, err
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	content, err := clientset.Discovery().RESTClient().Get().Timeout(5 * time.Second).AbsPath("/metrics").DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ScrapeKubeApiServerMetrics scrapes kube-apiserver /metrics endpoint and remote writes the metrics
func ScrapeKubeApiServerMetrics() error {
	kApiserverResp, err := ProbeKubeApiServerMetrics()
	if err != nil {
		return err
	}

	kApiserverMetrics, err := parseMetrics(kApiserverResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(kApiserverMetrics)

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
