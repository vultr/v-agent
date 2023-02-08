// Package config ensures the app is configured properly
package config

import (
	"errors"
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
