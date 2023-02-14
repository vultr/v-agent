// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// MemInfo is container for memory metrics
type MemInfo struct {
	MemTotal  uint64
	MemFree   uint64
	Buffers   uint64
	Cached    uint64
	SwapTotal uint64
	SwapFree  uint64
}

// getMeminfo reads /proc/meminfo and calculates the systems current memory
// statistics including swap
func getMeminfo() (*MemInfo, error) {
	_, err := os.Stat("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	procMemInfoFD, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer procMemInfoFD.Close() //nolint

	reader := bufio.NewReader(procMemInfoFD)

	var meminfo MemInfo

	for {
		data, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		splitted := strings.Fields(data)

		for _ = range splitted { //nolint
			if splitted[0] == "MemTotal:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.MemTotal = val * 1024
			} else if splitted[0] == "Buffers:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.Buffers = val * 1024
			} else if splitted[0] == "Cached:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.Cached = val * 1024
			} else if splitted[0] == "MemFree:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.MemFree = val * 1024
			} else if splitted[0] == "SwapTotal:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.SwapTotal = val * 1024
			} else if splitted[0] == "SwapFree:" {
				val, err1 := strconv.ParseUint(splitted[1], 10, 64)

				if err1 != nil {
					return nil, err1
				}

				meminfo.SwapFree = val * 1024
			}
		}
	}

	meminfo.MemFree = meminfo.MemFree + meminfo.Buffers + meminfo.Cached

	return &meminfo, nil
}
