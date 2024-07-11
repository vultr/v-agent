// Package metrics metrics collection
package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Loadavg from /proc/loadavg
type Loadavg struct {
	Load1        float64
	Load5        float64
	Load15       float64
	TasksRunning uint64
	TasksTotal   uint64
}

// getLoadavg returns the system load
func getLoadavg() (*Loadavg, error) {
	_, err := os.Stat("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	loadavgFile, err := os.Open("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	defer loadavgFile.Close() //nolint

	reader := bufio.NewReader(loadavgFile)

	data, err := reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	splitData := strings.Split(data, " ")

	if len(splitData) != 5 { //nolint
		return nil, fmt.Errorf("expect /proc/loadavg size to be length 5")
	}

	var ld Loadavg

	load1, err := strconv.ParseFloat(splitData[0], 64)
	if err != nil {
		return nil, err
	}

	load5, err := strconv.ParseFloat(splitData[1], 64)
	if err != nil {
		return nil, err
	}

	load15, err := strconv.ParseFloat(splitData[2], 64)
	if err != nil {
		return nil, err
	}

	tasks := strings.Split(splitData[3], "/")

	if len(tasks) != 2 { //nolint
		return nil, fmt.Errorf("expect /proc/loadavg tasks to be length 2")
	}

	running, err := strconv.ParseUint(tasks[0], 10, 64)
	if err != nil {
		return nil, err
	}

	total, err := strconv.ParseUint(tasks[1], 10, 64)
	if err != nil {
		return nil, err
	}

	ld.Load1 = load1
	ld.Load5 = load5
	ld.Load15 = load15
	ld.TasksRunning = running
	ld.TasksTotal = total

	return &ld, nil
}
