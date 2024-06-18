// Package metrics metrics collection
package metrics

import (
	"testing"
)

func TestGetLoadavg(t *testing.T) {
	_, err := getLoadavg()
	if err != nil {
		t.Error(err)
	}
}
