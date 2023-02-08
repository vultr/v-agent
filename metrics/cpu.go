// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ProcStatCPU container for CPU metrics
type ProcStatCPU struct {
	User, Nice, System, Idle, IOWait, IRQ, SoftIRQ, Steal, Guest, GuestNice int
}

func getCPUUtil() (*ProcStatCPU, error) {
	cpuStat1, err := getProcStat()
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second) // wait 1 second to capture another snapshot

	cpuStat2, err := getProcStat()
	if err != nil {
		return nil, err
	}

	// calculate delta between cpuStat2/cpuStat1
	deltaCPUStat := ProcStatCPU{
		User:      cpuStat2.User - cpuStat1.User,
		Nice:      cpuStat2.Nice - cpuStat1.Nice,
		System:    cpuStat2.System - cpuStat1.System,
		Idle:      cpuStat2.Idle - cpuStat1.Idle,
		IOWait:    cpuStat2.IOWait - cpuStat1.IOWait,
		IRQ:       cpuStat2.IRQ - cpuStat1.IRQ,
		SoftIRQ:   cpuStat2.SoftIRQ - cpuStat1.SoftIRQ,
		Steal:     cpuStat2.Steal - cpuStat1.Steal,
		Guest:     cpuStat2.Guest - cpuStat1.Guest,
		GuestNice: cpuStat2.GuestNice - cpuStat1.GuestNice,
	}

	return &deltaCPUStat, nil
}

func getProcStat() (*ProcStatCPU, error) {
	procStat, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer procStat.Close()

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

	if cpuStats.User, errAtoi = strconv.Atoi(splitted[1]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.Nice, errAtoi = strconv.Atoi(splitted[2]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.System, errAtoi = strconv.Atoi(splitted[3]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.Idle, errAtoi = strconv.Atoi(splitted[4]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.IOWait, errAtoi = strconv.Atoi(splitted[5]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.IRQ, errAtoi = strconv.Atoi(splitted[6]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.SoftIRQ, errAtoi = strconv.Atoi(splitted[7]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.Steal, errAtoi = strconv.Atoi(splitted[8]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.Guest, errAtoi = strconv.Atoi(splitted[9]); err != nil {
		return nil, errAtoi
	}

	if cpuStats.GuestNice, errAtoi = strconv.Atoi(splitted[10]); err != nil {
		return nil, errAtoi
	}

	return &cpuStats, nil
}

func getHostCPUs() uint {
	return uint(runtime.NumCPU())
}
