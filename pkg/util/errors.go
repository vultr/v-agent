// Package util provides utility functionality
package util

import "errors"

var (
	// ErrUnableToGetSUBID returned if unable to get subid
	ErrUnableToGetSUBID = errors.New("unable to probe subid")
)
