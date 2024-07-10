// Package metrics metrics collection
package metrics

import (
	"testing"
)

func TestGetMeminfo(t *testing.T) {
	_, err := getMeminfo()
	if err != nil {
		t.Error(err)
	}
}
