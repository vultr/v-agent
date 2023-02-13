// Package metrics provides prometheus metrics
package metrics

import (
	"context"
	"errors"
	"fmt"

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
func ProbeKubeApiServerMetrics() error {
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

	content, err := clientset.Discovery().RESTClient().Get().AbsPath("/metrics").DoRaw(context.TODO())
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", content)

	// TODO: Parse content and turn the metrics into something to send with remote write

	return nil
}
