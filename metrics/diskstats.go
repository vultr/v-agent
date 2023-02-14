// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"go.uber.org/zap"
)

var (
	// ErrInconsistentDisksReturned thrown if a read between /proc/diskstats results in less devices
	ErrInconsistentDisksReturned = errors.New("inconsistent amount of disks returned")
)

// DiskStats from /proc/diskstats
//
// https://www.kernel.org/doc/Documentation/iostats.txt
type DiskStats struct {
	Device                 string
	Reads                  uint64
	ReadsMerged            uint64
	SectorsRead            uint64
	MillisecondsReading    uint64
	WritesCompleted        uint64
	WritesMerged           uint64
	SectorsWritten         uint64
	MillisecondsWriting    uint64
	IOsInProgress          uint64
	MillisecondsInIOs      uint64
	WeightedIOsInMS        uint64
	Discards               uint64
	DiscardsMerged         uint64
	SectorsDiscarded       uint64
	MillisecondsDiscarding uint64
}

func getDiskStatsUtil() ([]*DiskStats, error) {
	diskStat1, err := getDiskStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second) // wait 1 second to capture another snapshot

	diskStat2, err := getDiskStats()
	if err != nil {
		return nil, err
	}

	var ds []*DiskStats

	if len(diskStat2) == len(diskStat1) {
		for i := range diskStat2 {
			if diskStat2[i].Device == diskStat1[i].Device {
				ds = append(ds, &DiskStats{
					Device:                 diskStat2[i].Device,
					Reads:                  diskStat2[i].Reads - diskStat1[i].Reads,
					ReadsMerged:            diskStat2[i].ReadsMerged - diskStat1[i].ReadsMerged,
					SectorsRead:            diskStat2[i].SectorsRead - diskStat1[i].SectorsRead,
					MillisecondsReading:    diskStat2[i].MillisecondsReading - diskStat1[i].MillisecondsReading,
					WritesCompleted:        diskStat2[i].WritesCompleted - diskStat1[i].WritesCompleted,
					WritesMerged:           diskStat2[i].WritesMerged - diskStat1[i].WritesMerged,
					SectorsWritten:         diskStat2[i].SectorsWritten - diskStat1[i].SectorsWritten,
					MillisecondsWriting:    diskStat2[i].MillisecondsWriting - diskStat1[i].MillisecondsWriting,
					IOsInProgress:          diskStat2[i].IOsInProgress,
					MillisecondsInIOs:      diskStat2[i].MillisecondsInIOs - diskStat1[i].MillisecondsInIOs,
					WeightedIOsInMS:        diskStat2[i].WeightedIOsInMS - diskStat1[i].WeightedIOsInMS,
					Discards:               diskStat2[i].Discards - diskStat1[i].Discards,
					DiscardsMerged:         diskStat2[i].DiscardsMerged - diskStat1[i].DiscardsMerged,
					SectorsDiscarded:       diskStat2[i].SectorsDiscarded - diskStat1[i].SectorsDiscarded,
					MillisecondsDiscarding: diskStat2[i].MillisecondsDiscarding - diskStat1[i].MillisecondsDiscarding,
				})
			} else {
				continue
			}
		}
	} else {
		return nil, ErrInconsistentDisksReturned
	}

	return ds, nil
}

// getDiskStats reads /proc/diskstats
func getDiskStats() ([]*DiskStats, error) {
	log := zap.L().Sugar()

	_, err := os.Stat("/proc/diskstats")
	if err != nil {
		return nil, err
	}

	procDiskStatsFD, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer procDiskStatsFD.Close() //nolint

	reader := bufio.NewReader(procDiskStatsFD)

	var diskStats []*DiskStats

	filter, err := config.GetDiskStatsFilter()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(*filter)

	for {
		data, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		splitted := strings.Fields(data)
		if len(splitted) != 18 { //nolint
			break
		}

		device := splitted[2]

		if re.MatchString(device) {
			log.Infof("skipping device %q", device)
			continue
		}

		readsS := splitted[3]
		readsSM := splitted[4]
		sectorsRS := splitted[5]
		msRdS := splitted[6]
		writesCompletedS := splitted[7]
		writesMergedS := splitted[8]
		sectorsWrS := splitted[9]
		msWrS := splitted[10]
		ioInProgressS := splitted[11]
		msInIOsS := splitted[12]
		weightedIOsInMSS := splitted[13]
		discardsS := splitted[14]
		discardsMergedS := splitted[15]
		sectorsDiscardedS := splitted[16]
		msDiscardingS := splitted[17]

		reads, err := strconv.ParseUint(readsS, 10, 64)
		if err != nil {
			return nil, err
		}
		readsM, err := strconv.ParseUint(readsSM, 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsRd, err := strconv.ParseUint(sectorsRS, 10, 64)
		if err != nil {
			return nil, err
		}
		msRd, err := strconv.ParseUint(msRdS, 10, 64)
		if err != nil {
			return nil, err
		}
		writesCompleted, err := strconv.ParseUint(writesCompletedS, 10, 64)
		if err != nil {
			return nil, err
		}
		writesMerged, err := strconv.ParseUint(writesMergedS, 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsWr, err := strconv.ParseUint(sectorsWrS, 10, 64)
		if err != nil {
			return nil, err
		}
		msWr, err := strconv.ParseUint(msWrS, 10, 64)
		if err != nil {
			return nil, err
		}
		ioInProgress, err := strconv.ParseUint(ioInProgressS, 10, 64)
		if err != nil {
			return nil, err
		}
		msInIOs, err := strconv.ParseUint(msInIOsS, 10, 64)
		if err != nil {
			return nil, err
		}
		weightedIOsInMS, err := strconv.ParseUint(weightedIOsInMSS, 10, 64)
		if err != nil {
			return nil, err
		}
		discards, err := strconv.ParseUint(discardsS, 10, 64)
		if err != nil {
			return nil, err
		}
		discardsMerged, err := strconv.ParseUint(discardsMergedS, 10, 64)
		if err != nil {
			return nil, err
		}
		sectorsDiscarded, err := strconv.ParseUint(sectorsDiscardedS, 10, 64)
		if err != nil {
			return nil, err
		}
		msDiscarding, err := strconv.ParseUint(msDiscardingS, 10, 64)
		if err != nil {
			return nil, err
		}

		var ds DiskStats

		ds.Device = device
		ds.Reads = reads
		ds.ReadsMerged = readsM
		ds.SectorsRead = sectorsRd
		ds.MillisecondsReading = msRd
		ds.WritesCompleted = writesCompleted
		ds.WritesMerged = writesMerged
		ds.SectorsWritten = sectorsWr
		ds.MillisecondsWriting = msWr
		ds.IOsInProgress = ioInProgress
		ds.MillisecondsInIOs = msInIOs
		ds.WeightedIOsInMS = weightedIOsInMS
		ds.Discards = discards
		ds.DiscardsMerged = discardsMerged
		ds.SectorsDiscarded = sectorsDiscarded
		ds.MillisecondsDiscarding = msDiscarding

		diskStats = append(diskStats, &ds)
	}

	return diskStats, nil
}
