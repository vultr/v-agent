// Package config ensures the app is configured properly
package config

import (
	"github.com/vultr/v-agent/pkg/util"

	"go.uber.org/zap"
)

// GetConfig returns config
func GetConfig() *Config {
	return &cfg
}

// GetVersion returns application version
func GetVersion() string {
	cfg := GetConfig()

	return cfg.Version
}

// GetLabels returns all labels
func GetLabels() map[string]string {
	return labels
}

// GetLabel returns a map[string]string of the requested label or error
func GetLabel(label string) map[string]string {
	log := zap.L().Sugar()

	l := make(map[string]string)
	for k, v := range labels {
		if k == label {
			l[k] = v
			return l
		}
	}

	log.Warn("label %s does not exist", label)

	l[label] = "unknown"

	return l
}

// GetProbesAPIListen returns probes api listen addr
func GetProbesAPIListen() string {
	cfg := GetConfig()

	return cfg.ProbesAPI.Listen
}

// GetProbesAPIPort returns probes api listen port
func GetProbesAPIPort() uint {
	cfg := GetConfig()

	return cfg.ProbesAPI.Port
}

// GetDiskStatsFilter returns the regex for the disk stats filter
func GetDiskStatsFilter() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.DiskStats.Filter
}

// GetKubeconfig returns path to kubeconfig
func GetKubeconfig() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Kubernetes.Kubeconfig
}

// GetKonnectivityHealthEndpoint returns health endpoint
func GetKonnectivityHealthEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Konnectivity.HealthEndpoint
}

// GetKonnectivityMetricsEndpoint returns metrics endpoint
func GetKonnectivityMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Konnectivity.MetricsEndpoint
}

// GetEtcdEndpoint returns etcd endpoint
func GetEtcdEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Etcd.Endpoint
}

// GetEtcdCACert returns path to cacert
func GetEtcdCACert() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Etcd.CACert
}

// GetEtcdClientCert returns path to client cert
func GetEtcdClientCert() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Etcd.Cert
}

// GetEtcdClientKey returns path to client key
func GetEtcdClientKey() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Etcd.Key
}

// LoadAvgMetricCollectionEnabled returns true/false if load_avg metrics collection enabled
func LoadAvgMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.LoadAvg.Enabled
}

// CPUMetricCollectionEnabled returns true/false if cpu metrics collection enabled
func CPUMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.CPU.Enabled
}

// MemoryMetricCollectionEnabled returns true/false if memory metrics collection enabled
func MemoryMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Memory.Enabled
}

// NICMetricCollectionEnabled returns true/false if nic collection enabled
func NICMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.NIC.Enabled
}

// DiskStatsMetricCollectionEnabled returns true/false if diskstats collection enabled
func DiskStatsMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.DiskStats.Enabled
}

// FileSystemMetricCollectionEnabled returns true/false if filesystem collection enabled
func FileSystemMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Filesystem.Enabled
}

// KubernetesMetricCollectionEnabled returns true/false if Kubernetes collection enabled
func KubernetesMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Kubernetes.Enabled
}

// KonnectivityMetricCollectionEnabled returns true/false if Konnectivity collection enabled
func KonnectivityMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Konnectivity.Enabled
}

// EtcdMetricCollectionEnabled returns true/false if Etcd collection enabled
func EtcdMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Etcd.Enabled
}

// NginxVTSMetricsCollectionEnabled returns true if nginx vts metric collection enabled
func NginxVTSMetricsCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.NginxVTS.Enabled
}

// GetNginxVTSMetricsEndpoint returns nginx vts endpoint
func GetNginxVTSMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.NginxVTS.Endpoint
}

// VCDNAgentMetricsCollectionEnabled returns true if v-cdn-agent metric collection enabled
func VCDNAgentMetricsCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.VCDNAgent.Enabled
}

// GetVCDNAgentMetricsEndpoint returns v-cdn-agent endpoint
func GetVCDNAgentMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.VCDNAgent.Endpoint
}

// HAProxyMetricCollectionEnabled returns true/false if haproxy collection enabled
func HAProxyMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.HAProxy.Enabled
}

// GetHAProxyMetricsEndpoint returns haproxy endpoint
func GetHAProxyMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.HAProxy.Endpoint
}

// CephMetricCollectionEnabled returns true/false if ceph collection enabled
func CephMetricCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Ceph.Enabled
}

// GetCephMetricsEndpoint returns metrics endpoint
func GetCephMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.Ceph.Endpoint
}

// VDNSMetricsCollectionEnabled returns true if v-dns metric collection enabled
func VDNSMetricsCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.VDNS.Enabled
}

// GetVDNSMetricsEndpoint returns v-dns endpoint
func GetVDNSMetricsEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.VDNS.Endpoint
}

// KubernetesPodsCollectionEnabled returns true if pod scrape collection is enabled
func KubernetesPodsCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Kubernetes.Pods.Enabled
}

// GetKubernetesPodsNamespaces returns namespaces to scrape pods from
func GetKubernetesPodsNamespaces() []string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Kubernetes.Pods.Namespaces
}

// SMARTCollectionEnabled returns true if SMART collection is enabled
func SMARTCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Agent.SMART.Enabled
}

// GetSMARTBlockDevices returns SMART block devices to scan
func GetSMARTBlockDevices() []string {
	log := zap.L().Sugar()

	cfg := GetConfig()

	// block_devices provided, use those instead
	if len(cfg.MetricsConfig.Agent.SMART.BlockDevices) > 0 {
		return cfg.MetricsConfig.Agent.SMART.BlockDevices
	}

	// block_devices not provided, scan and build list
	blockDevices, err := util.GenerateBlockDevices()
	if err != nil {
		log.Warn(err)

		return []string{}
	}

	return blockDevices
}

// DCGMCollectionEnabled returns true if DCGM collection is enabled
func DCGMCollectionEnabled() bool {
	cfg := GetConfig()

	return cfg.MetricsConfig.Kubernetes.DCGM.Enabled
}

// GetDCGMNamespace returns the DCGM namespace
func GetDCGMNamespace() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Kubernetes.DCGM.Namespace
}

// GetDCGMEndpoint returns the DCGM endpoint
func GetDCGMEndpoint() string {
	cfg := GetConfig()

	return cfg.MetricsConfig.Kubernetes.DCGM.Endpoint
}
