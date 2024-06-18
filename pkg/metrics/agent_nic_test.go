// Package metrics metrics collection
package metrics

import (
	"testing"
)

func TestGetNICStats(t *testing.T) {
	_, err := getNICStats()
	if err != nil {
		t.Error(err)
	}
}
