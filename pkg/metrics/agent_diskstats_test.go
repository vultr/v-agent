// Package metrics metrics collection
package metrics

import (
	"reflect"
	"testing"

	"github.com/vultr/v-agent/cmd/v-agent/config"
)

func TestGetDiskStatsUtil(t *testing.T) {
	config.NewConfig("test", "v0.0.0") //nolint

	p, err := getDiskStatsUtil()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(p, prevDS) {
		t.Errorf("expect prevDS and p to be equal")
	}
}

func TestGetDiskStats(t *testing.T) {
	_, err := getDiskStats()
	if err != nil {
		t.Error(err)
	}
}
