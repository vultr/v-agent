// Package metrics metrics collection
package metrics

import "errors"

var (
	// ErrInconsistentDisksReturned thrown if a read between /proc/diskstats results in less devices
	ErrInconsistentDisksReturned = errors.New("inconsistent amount of disks returned")

	// ErrCephMgrNotActive thrown if no data returned when attempting to read metrics
	ErrCephMgrNotActive = errors.New("ceph-mgr not active")

	// ErrEtcdUnhealthy returned if response is not "true" from /health
	ErrEtcdUnhealthy = errors.New("etcd unhealthy")

	// ErrHAProxyServerUnhealthy returned if response is not status code 200 from /metrics
	ErrHAProxyServerUnhealthy = errors.New("haproxy unhealthy")

	// ErrKonnectivityServerUnhealthy returned if response is not "ok" from /healthz
	ErrKonnectivityServerUnhealthy = errors.New("konnectivity unhealthy")

	// ErrKubeAPIServerUnhealthy returned if response is not "ok" from /healthz
	ErrKubeAPIServerUnhealthy = errors.New("kube-apiserver unhealthy")

	// ErrNginxVTSServerUnhealthy returned if response is not status code 200 from /metrics
	ErrNginxVTSServerUnhealthy = errors.New("nginx-vts unhealthy")

	// ErrVCDNAgentServerUnhealthy returned if response is not status code 200 from /metrics
	ErrVCDNAgentServerUnhealthy = errors.New("v-cdn-agent unhealthy")

	// ErrVDNSUnhealthy returned if response is not status code 200 from /metrics
	ErrVDNSUnhealthy = errors.New("v-dns unhealthy")
)
