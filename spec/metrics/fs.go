// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"io"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

// Mounts from /proc/mounts
type Mounts struct {
	Device    string
	Mount     string
	Type      string
	Opts      string
	Dump      bool
	BootCheck bool
}

// FilesystemStats stats
type FilesystemStats struct {
	Device     string
	Mount      string
	BlockSize  int64
	Inodes     uint64
	InodesFree uint64
	InodesUsed uint64
	InodesUtil float64
	BytesTotal uint64
	BytesFree  uint64
	BytesUsed  uint64
	BytesUtil  float64
}

// getFilesystemUtil calls statfs syscall to get utilization
func getFilesystemUtil() ([]*FilesystemStats, error) {
	log := zap.L().Sugar()

	mounts, err := getMounts()
	if err != nil {
		return nil, err
	}

	var fs []*FilesystemStats

	for i := range mounts {
		df := syscall.Statfs_t{}

		if err := syscall.Statfs(mounts[i].Mount, &df); err != nil {
			log.Error(err)
			continue
		}

		totalBlocks := df.Blocks
		freeBlocks := df.Bfree
		bytesUsed := (uint64(df.Bsize) * totalBlocks) - (uint64(df.Bsize) * freeBlocks)

		fs = append(fs, &FilesystemStats{
			Device:     mounts[i].Device,
			Mount:      mounts[i].Mount,
			BlockSize:  df.Bsize,
			Inodes:     df.Files,
			InodesFree: df.Ffree,
			InodesUsed: df.Files - df.Ffree,
			InodesUtil: ((float64(df.Files) - float64(df.Ffree)) / float64(df.Files)) * 100,
			BytesTotal: uint64(df.Bsize) * totalBlocks,
			BytesFree:  uint64(df.Bsize) * freeBlocks,
			BytesUsed:  bytesUsed,
			BytesUtil:  float64(bytesUsed) / (float64(df.Bsize) * float64(totalBlocks)) * 100,
		})
	}

	return fs, nil
}

// getMounts reads /proc/mounts
func getMounts() ([]*Mounts, error) {
	_, err := os.Stat("/proc/mounts")
	if err != nil {
		return nil, err
	}

	mountsFD, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer mountsFD.Close() //nolint

	reader := bufio.NewReader(mountsFD)

	var mounts []*Mounts

	for {
		data, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		splitted := strings.Fields(data)
		if len(splitted) != 6 { //nolint
			break
		}

		device := splitted[0]
		mount := splitted[1]
		typefs := splitted[2]
		opts := splitted[3]

		var dump, boot bool
		if splitted[4] == "0" {
			dump = false
		} else if splitted[4] == "1" {
			dump = true
		}

		if splitted[5] == "0" {
			boot = false
		} else if splitted[5] == "1" {
			boot = true
		}

		if strings.HasPrefix(device, "/") {
			mounts = append(mounts, &Mounts{
				Device:    device,
				Mount:     mount,
				Type:      typefs,
				Opts:      opts,
				Dump:      dump,
				BootCheck: boot,
			})
		}
	}

	return mounts, nil
}
