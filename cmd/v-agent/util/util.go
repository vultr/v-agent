// Package util provides utility functionality
package util

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

var (
	// ErrUnableToGetSUBID returned if unable to get subid
	ErrUnableToGetSUBID = errors.New("unable to probe subid")
)

// GetVPSID returns the VPS ID from http://169.254.169.254/v1.json .instanceid
func GetVPSID() (*string, error) {
	resp, err := http.Get("http://169.254.169.254/v1.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	vpsid := gjson.Get(string(body), "instanceid")

	return &vpsid.Str, nil
}

// GetSubID returns the subid of the underlying service
//
// vke has one method of extraction of the subid
// vlb has another method of extraction of the subid
// vfs will most likely have its own method of extraction of the subid
func GetSubID(product string) (*string, error) {
	switch product {
	case "vke":
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
	case "vlb":
		resp, err := http.Get("http://169.254.169.254/latest/user-data")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close() //nolint

		realURL, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		resp2, err := http.Get(string(realURL))
		if err != nil {
			return nil, err
		}
		defer resp2.Body.Close() //nolint

		body, err := io.ReadAll(resp2.Body)
		if err != nil {
			return nil, err
		}

		subid := gjson.Get(string(body), "load_balancer_config.lb_subid")

		return &subid.Raw, nil
	case "vfs":
		resp, err := http.Get("http://169.254.169.254/latest/user-data")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close() //nolint

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// TODO needs implementation

		subid := gjson.Get(string(body), "data.vfs_subid")

		return &subid.Str, nil
	}

	return nil, ErrUnableToGetSUBID
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
