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

// NICMetricCollectionEnabled returns true/false if memory nic collection enabled
func NICMetricCollectionEnabled() bool {
	log := zap.L().Sugar()

	cfg, err := GetConfig()
	if err != nil {
		log.Error(err)
		return true
	}

	return cfg.MetricsConfig.NIC.Enabled
}
