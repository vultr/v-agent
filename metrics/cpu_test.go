// Package metrics provides prometheus metrics
package metrics

import (
	"reflect"
	"testing"
)

func TestGetCPUUtil(t *testing.T) {
	p, err := getCPUUtil()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*p, prevCPU) {
		t.Errorf("expect prevCPU and p to be equal")
	}
}

func TestGetProcStat(t *testing.T) {
	_, err := getProcStat()
	if err != nil {
		t.Error(err)
	}
}

func TestGetHostCPUs(t *testing.T) {
	cpus := getHostCPUs()
	if cpus < 1 {
		t.Error("expect at least 1 cpu")
	}
}
