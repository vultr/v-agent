// Package config ensures the app is configured properly
package config

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	c, _ := GetConfig()
	if c != nil {
		t.Error("expecting GetConfig to return nil")
	}

	NewConfig("test", "v0.0.0") //nolint
	c1, _ := GetConfig()
	if c1 == nil {
		t.Error("expecting GetConfig to not be nil")
	}
}
