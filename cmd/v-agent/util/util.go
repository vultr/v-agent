// Package util provides utility functionality
package util

import (
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// GetSubID returns the SUBID extracted from http://169.254.169.254/latest/user-data | jq '.data.vke_subid'
func GetSubID() (*string, error) {
	resp, err := http.Get("http://169.254.169.254/latest/user-data")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	subid := gjson.Get(string(body), "data.vke_subid")

	return &subid.Str, nil
}
