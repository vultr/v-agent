// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/vultr/v-agent/cmd/v-agent/config"

	"go.uber.org/zap"
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

var prevDS []*DiskStats

func getDiskStatsUtil() ([]*DiskStats, error) {
	diskStat1, err := getDiskStats()
	if err != nil {
		return nil, err
	}

	var ds []*DiskStats

	for i := range diskStat1 {
		for j := range prevDS {
			if diskStat1[i].Device == prevDS[j].Device {
				ds = append(ds, &DiskStats{
					Device:                 diskStat1[i].Device,
					Reads:                  diskStat1[i].Reads - prevDS[j].Reads,
					ReadsMerged:            diskStat1[i].ReadsMerged - prevDS[j].ReadsMerged,
					SectorsRead:            diskStat1[i].SectorsRead - prevDS[j].SectorsRead,
					MillisecondsReading:    diskStat1[i].MillisecondsReading - prevDS[j].MillisecondsReading,
					WritesCompleted:        diskStat1[i].WritesCompleted - prevDS[j].WritesCompleted,
					WritesMerged:           diskStat1[i].WritesMerged - prevDS[j].WritesMerged,
					SectorsWritten:         diskStat1[i].SectorsWritten - prevDS[j].SectorsWritten,
					MillisecondsWriting:    diskStat1[i].MillisecondsWriting - prevDS[j].MillisecondsWriting,
					IOsInProgress:          diskStat1[i].IOsInProgress,
					MillisecondsInIOs:      diskStat1[i].MillisecondsInIOs - prevDS[j].MillisecondsInIOs,
					WeightedIOsInMS:        diskStat1[i].WeightedIOsInMS - prevDS[j].WeightedIOsInMS,
					Discards:               diskStat1[i].Discards - prevDS[j].Discards,
					DiscardsMerged:         diskStat1[i].DiscardsMerged - prevDS[j].DiscardsMerged,
					SectorsDiscarded:       diskStat1[i].SectorsDiscarded - prevDS[j].SectorsDiscarded,
					MillisecondsDiscarding: diskStat1[i].MillisecondsDiscarding - prevDS[j].MillisecondsDiscarding,
				})
			}
		}
	}

	prevDS = diskStat1

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

	filter := config.GetDiskStatsFilter()

	re := regexp.MustCompile(filter)

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
