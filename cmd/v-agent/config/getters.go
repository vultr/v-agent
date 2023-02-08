// Package config ensures the app is configured properly
package config

import (
	"errors"
)

var (
	ErrConfigNotInitialized = errors.New("config not initialized")
)

// GetConfig returns config
func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	return nil, ErrConfigNotInitialized
}
