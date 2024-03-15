// Package config ensures the app is configured properly
package config

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	NewConfig("test", "v0.0.0") //nolint
	c1 := GetConfig()
	if c1 == nil {
		t.Error("expecting GetConfig to not be nil")
	}
}
