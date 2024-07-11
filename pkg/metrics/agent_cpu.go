// Package metrics metrics collection
package metrics

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// ProcStatCPU container for CPU metrics
type ProcStatCPU struct {
	User, Nice, System, Idle, IOWait, IRQ, SoftIRQ, Steal, Guest, GuestNice int
}

var prevCPU ProcStatCPU

func getCPUUtil() (*ProcStatCPU, error) {
	cpuStat1, err := getProcStat()
	if err != nil {
		return nil, err
	}

	stat := &ProcStatCPU{
		User:      cpuStat1.User - prevCPU.User,
		Nice:      cpuStat1.Nice - prevCPU.Nice,
		System:    cpuStat1.System - prevCPU.System,
		Idle:      cpuStat1.Idle - prevCPU.Idle,
		IOWait:    cpuStat1.IOWait - prevCPU.IOWait,
		IRQ:       cpuStat1.IRQ - prevCPU.IRQ,
		SoftIRQ:   cpuStat1.SoftIRQ - prevCPU.SoftIRQ,
		Steal:     cpuStat1.Steal - prevCPU.Steal,
		Guest:     cpuStat1.Guest - prevCPU.Guest,
		GuestNice: cpuStat1.GuestNice - prevCPU.GuestNice,
	}

	prevCPU = *cpuStat1

	return stat, nil
}

func getProcStat() (*ProcStatCPU, error) {
	procStat, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer procStat.Close() //nolint

	reader := bufio.NewReader(procStat)

	data, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	splitted := strings.Fields(data)

	if splitted[0] != "cpu" {
		return nil, fmt.Errorf("metrics: expecting first field to be 'cpu' when reading /proc/stat")
	}

	if len(splitted) != 11 { //nolint
		return nil, fmt.Errorf("metrics: expecting 'cpu' from /proc/stat to be 11 (%d)", len(splitted))
	}

	var cpuStats ProcStatCPU
	var errAtoi error

	if cpuStats.User, errAtoi = strconv.Atoi(splitted[1]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.Nice, errAtoi = strconv.Atoi(splitted[2]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.System, errAtoi = strconv.Atoi(splitted[3]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.Idle, errAtoi = strconv.Atoi(splitted[4]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.IOWait, errAtoi = strconv.Atoi(splitted[5]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.IRQ, errAtoi = strconv.Atoi(splitted[6]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.SoftIRQ, errAtoi = strconv.Atoi(splitted[7]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.Steal, errAtoi = strconv.Atoi(splitted[8]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.Guest, errAtoi = strconv.Atoi(splitted[9]); errAtoi != nil {
		return nil, errAtoi
	}

	if cpuStats.GuestNice, errAtoi = strconv.Atoi(splitted[10]); errAtoi != nil {
		return nil, errAtoi
	}

	return &cpuStats, nil
}

func getHostCPUs() uint {
	return uint(runtime.NumCPU())
}
