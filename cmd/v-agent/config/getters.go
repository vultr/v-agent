// Package config ensures the app is configured properly
package config

import (
	"errors"

	"go.uber.org/zap"
)

var (
	// ErrConfigNotInitialized returned if the configuration is not initialized
	ErrConfigNotInitialized = errors.New("config not initialized")
	// ErrSubIDNotSet returned if the subid is empty
	ErrSubIDNotSet = errors.New("subid is not set")
)

// GetConfig returns config
func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	return nil, ErrConfigNotInitialized
}

// GetSubID returns the subid
func GetSubID() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	if cfg.SubID == "" {
		return nil, ErrSubIDNotSet
	}

	return &cfg.SubID, nil
}

// GetProduct returns underlying product name
func GetProduct() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.Product, nil
}

// GetDiskStatsFilter returns the regex for the disk stats filter
func GetDiskStatsFilter() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.DiskStats.Filter, nil
}

// GetKubeconfig returns path to kubeconfig
func GetKubeconfig() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.Kubernetes.Kubeconfig, nil
}

// GetEtcdEndpoint returns etcd endpoint
func GetEtcdEndpoint() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.Etcd.Endpoint, nil
}

// GetEtcdCACert returns path to cacert
func GetEtcdCACert() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.Etcd.CACert, nil
}

// GetEtcdClientCert returns path to client cert
func GetEtcdClientCert() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.Etcd.Cert, nil
}

// GetEtcdClientKey returns path to client key
func GetEtcdClientKey() (*string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &cfg.MetricsConfig.Etcd.Key, nil
}

// LoadAvgMetricCollectionEnabled returns true/false if load_avg metrics collection enabled
func LoadAvgMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.LoadAvg.Enabled
}

// CPUMetricCollectionEnabled returns true/false if cpu metrics collection enabled
func CPUMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.CPU.Enabled
}

// MemoryMetricCollectionEnabled returns true/false if memory metrics collection enabled
func MemoryMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.Memory.Enabled
}

// NICMetricCollectionEnabled returns true/false if nic collection enabled
func NICMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.NIC.Enabled
}

// DiskStatsMetricCollectionEnabled returns true/false if diskstats collection enabled
func DiskStatsMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.DiskStats.Enabled
}

// FileSystemMetricCollectionEnabled returns true/false if filesystem collection enabled
func FileSystemMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.Filesystem.Enabled
}

// KubernetesMetricCollectionEnabled returns true/false if Kubernetes collection enabled
func KubernetesMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.Kubernetes.Enabled
}

// EtcdMetricCollectionEnabled returns true/false if Etcd collection enabled
func EtcdMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.Etcd.Enabled
}
