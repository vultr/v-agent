package connectors

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetKubernetesConn returns a connection to the kubernetes cluster
func GetKubernetesConn() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
