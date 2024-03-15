// Package util provides utility functionality
package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
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

		subid := gjson.Get(string(body), "data.vfs_subid")

		return &subid.Str, nil
	case "vcdn":
		resp, err := http.Get("http://169.254.169.254/latest/user-data")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close() //nolint

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

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

// GenerateBlockDevices generates block devices from /dev, not including any of the following: loop, dm, nbd, sr
func GenerateBlockDevices() ([]string, error) {
	log := zap.L().Sugar()

	files, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	var blockDevices []string
	for i := range files {
		if strings.HasPrefix(files[i].Name(), "dm") || strings.HasPrefix(files[i].Name(), "loop") || strings.HasPrefix(files[i].Name(), "nbd") || strings.HasPrefix(files[i].Name(), "sr") || strings.HasPrefix(files[i].Name(), "vd") {
			continue
		}

		absPath := fmt.Sprintf("/dev/%s", files[i].Name())

		if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
			log.Warnf("block device %s does not exist", absPath)

			continue
		}

		blockDevices = append(blockDevices, absPath)
	}

	return blockDevices, nil
}
