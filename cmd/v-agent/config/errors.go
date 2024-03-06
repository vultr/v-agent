// Package config ensures the app is configured properly
package config

import "errors"

var (
	// ErrConfigNotInitialized returned if the configuration is not initialized
	ErrConfigNotInitialized = errors.New("config not initialized")
	// ErrSubIDNotSet returned if the subid is empty
	ErrSubIDNotSet = errors.New("subid is not set")
	// ErrVPSIDNotSet returned if the vpsid is empty
	ErrVPSIDNotSet = errors.New("vpsid is not set")
	// ErrLabelNotExist returned when the specified label doesnt exist
	ErrLabelNotExist = errors.New("label does not exist")
)
