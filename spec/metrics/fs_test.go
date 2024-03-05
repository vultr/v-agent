// Package metrics provides prometheus metrics
package metrics

import (
	"testing"
)

func TestGetFilesystemUtil(t *testing.T) {
	_, err := getFilesystemUtil()
	if err != nil {
		t.Error(err)
	}
}

func TestGetMounts(t *testing.T) {
	_, err := getMounts()
	if err != nil {
		t.Error(err)
	}
}
