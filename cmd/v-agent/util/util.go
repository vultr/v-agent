// Package util provides utility functionality
package util

import (
	"bytes"
	"io"
	"net/http"
	"os"

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

// GetBIOSVendor returns bios vendor
func GetBIOSVendor() (*string, error) {
	vendor, err := os.ReadFile("/sys/devices/virtual/dmi/id/bios_vendor")
	if err != nil {
		return nil, err
	}

	vendor = bytes.Trim(vendor, "\x00") // NUL
	vendor = bytes.Trim(vendor, "\x0A") // \n
	vendor = bytes.Trim(vendor, "\x0D") // \r

	s := string(vendor)

	return &s, nil
}
