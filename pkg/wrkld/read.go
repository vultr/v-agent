package wrkld

import (
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetScrapeablePods returns all pods in the namespce that have prometheus.io/scrape=true
func GetScrapeablePods(client kubernetes.Interface, namespace string) ([]v1.Pod, error) {
	log := zap.L().Sugar()

	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var scrapeablePodsItems []v1.Pod
	for i := range pods.Items {
		vv1, ok1 := pods.Items[i].Annotations["prometheus.io/scrape"]
		if !ok1 {
			log.Infof("pod %s doesn't have prometheus.io/scrape annotation", pods.Items[i].ObjectMeta.Name)

			continue
		}

		if !strings.EqualFold(vv1, "true") {
			log.Infof("pod %s prometheus.io/scrape is not true", pods.Items[i].ObjectMeta.Name)

			continue
		}

		v2, ok3 := pods.Items[i].Annotations["prometheus.io/port"]
		if !ok3 {
			log.Infof("pod %s doesn't have prometheus.io/port annotation", pods.Items[i].ObjectMeta.Name)

			continue
		}

		_, err := strconv.Atoi(v2)
		if err != nil {
			log.Warnf("pod %s prometheus.io/port is not valid (%s)", pods.Items[i].ObjectMeta.Name, v2)

			continue
		}

		v3, ok2 := pods.Items[i].Annotations["prometheus.io/path"]
		if !ok2 {
			log.Infof("pod %s doesn't have prometheus.io/path annotation, default will be /metrics", pods.Items[i].ObjectMeta.Name)
		} else {
			if !strings.HasPrefix(v3, "/") {
				log.Warnf("pod %s prometheus.io/path is not valid (%s)", pods.Items[i].ObjectMeta.Name, v3)

				continue
			}
		}

		scrapeablePodsItems = append(scrapeablePodsItems, pods.Items[i])
	}

	return scrapeablePodsItems, nil
}

// GetEndpoint returns specified endpoint
func GetEndpoint(client kubernetes.Interface, namespace, name string) (*v1.Endpoints, error) {
	return client.CoreV1().Endpoints(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// GetNamespace returns specified namespace
func GetNamespace(client kubernetes.Interface, namespace string) (*v1.Namespace, error) {
	return client.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
}

// NamespaceExists returns true if the namespace exists
func NamespaceExists(client kubernetes.Interface, namespace string) bool {
	_, err := GetNamespace(client, namespace)

	return !errors.IsNotFound(err)
}

// EndpointExists returns true if the endpoint exists
func EndpointExists(client kubernetes.Interface, namespace, endpoint string) bool {
	_, err := GetEndpoint(client, namespace, endpoint)

	return !errors.IsNotFound(err)
}
