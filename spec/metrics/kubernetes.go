// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// DoKubeAPIServerHealthCheck probes /healthz and returns nil or ErrKubeAPIServerUnhealthy, or some other error
func DoKubeAPIServerHealthCheck() error {
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
	}

	return ErrKubeAPIServerUnhealthy
}

// ProbeKubeAPIServerMetrics probes /metrics from kube-apiserver
func ProbeKubeAPIServerMetrics() ([]byte, error) {
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

	content, err := clientset.Discovery().RESTClient().Get().Timeout(5 * time.Second).AbsPath("/metrics").DoRaw(context.TODO()) //nolint
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ScrapeKubeAPIServerMetrics scrapes kube-apiserver /metrics endpoint and remote writes the metrics
func ScrapeKubeAPIServerMetrics() error {
	apiserverResp, err := ProbeKubeAPIServerMetrics()
	if err != nil {
		return err
	}

	apiserverMetrics, err := parseMetrics(apiserverResp)
	if err != nil {
		return err
	}

	tsList := GetMetricsAsTimeSeries(apiserverMetrics)

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
