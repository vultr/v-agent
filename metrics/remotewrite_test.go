// Package metrics provides prometheus metrics
package metrics

import (
	"testing"
)

func TestValidateResponseStatus(t *testing.T) {
	err := validateResponseStatus(200)
	if err != nil {
		t.Error(err)
	}

	err2 := validateResponseStatus(500)
	if err2 == nil {
		t.Error("expect error when status code is not between 200 and 300")
	}
}
