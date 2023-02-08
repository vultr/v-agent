// Package metrics provides prometheus metrics
package metrics

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// NICMetrics metrics for the underlying NIC for the guest
type NICMetrics struct {
	Interface    string
	Bytes        uint
	BytesTX      uint
	BytesRX      uint
	Packets      uint
	PacketsTX    uint
	PacketsRX    uint
	Errors       uint
	ErrorsTX     uint
	ErrorsRX     uint
	Drop         uint
	DropTX       uint
	DropRX       uint
	FIFO         uint
	FIFOTX       uint
	FIFORX       uint
	FrameRX      uint
	CollsTX      uint
	Compressed   uint
	CompressedTX uint
	CompressedRX uint
	CarrierTX    uint
	MulticastRX  uint
}

func getNICStats() ([]NICMetrics, error) {
	procNetDev, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer procNetDev.Close()

	reader := bufio.NewReader(procNetDev)

	var metrics []NICMetrics

	for {
		data, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		} else if err == io.EOF {
			return metrics, nil
		}

		splitted := strings.Fields(data)

		if len(splitted) != 17 {
			continue
		}

		iface := splitted[0][:len(splitted[0])-1]

		var errAtoi error
		var bytesRX, bytesTX, packetsRX, packetsTX, errorsRX, errorsTX, dropRX, dropTX, fifoRX, fifoTX, frameRX, collsTX, compRX, compTX, multicastRX, carrierTX int
		if bytesRX, err = strconv.Atoi(splitted[1]); err != nil {
			return nil, errAtoi
		}
		if bytesTX, err = strconv.Atoi(splitted[9]); err != nil {
			return nil, errAtoi
		}
		if packetsRX, err = strconv.Atoi(splitted[2]); err != nil {
			return nil, errAtoi
		}
		if packetsTX, err = strconv.Atoi(splitted[10]); err != nil {
			return nil, errAtoi
		}
		if errorsRX, err = strconv.Atoi(splitted[3]); err != nil {
			return nil, errAtoi
		}
		if errorsTX, err = strconv.Atoi(splitted[11]); err != nil {
			return nil, errAtoi
		}
		if dropRX, err = strconv.Atoi(splitted[4]); err != nil {
			return nil, errAtoi
		}
		if dropTX, err = strconv.Atoi(splitted[12]); err != nil {
			return nil, errAtoi
		}
		if fifoRX, err = strconv.Atoi(splitted[5]); err != nil {
			return nil, errAtoi
		}
		if fifoTX, err = strconv.Atoi(splitted[13]); err != nil {
			return nil, errAtoi
		}
		if frameRX, err = strconv.Atoi(splitted[6]); err != nil {
			return nil, errAtoi
		}
		if collsTX, err = strconv.Atoi(splitted[14]); err != nil {
			return nil, errAtoi
		}
		if compRX, err = strconv.Atoi(splitted[7]); err != nil {
			return nil, errAtoi
		}
		if compTX, err = strconv.Atoi(splitted[16]); err != nil {
			return nil, errAtoi
		}
		if multicastRX, err = strconv.Atoi(splitted[8]); err != nil {
			return nil, errAtoi
		}
		if carrierTX, err = strconv.Atoi(splitted[15]); err != nil {
			return nil, errAtoi
		}

		var gnm NICMetrics
		gnm.Interface = iface
		gnm.BytesRX = uint(bytesRX)
		gnm.BytesTX = uint(bytesTX)
		gnm.Bytes = gnm.BytesRX + gnm.BytesTX
		gnm.PacketsRX = uint(packetsRX)
		gnm.PacketsTX = uint(packetsTX)
		gnm.Packets = gnm.PacketsRX + gnm.PacketsTX
		gnm.ErrorsRX = uint(errorsRX)
		gnm.ErrorsTX = uint(errorsTX)
		gnm.Errors = gnm.ErrorsRX + gnm.ErrorsTX
		gnm.DropRX = uint(dropRX)
		gnm.DropTX = uint(dropTX)
		gnm.Drop = gnm.DropRX + gnm.DropTX
		gnm.FIFORX = uint(fifoRX)
		gnm.FIFOTX = uint(fifoTX)
		gnm.FIFO = gnm.FIFORX + gnm.FIFOTX
		gnm.FrameRX = uint(frameRX)
		gnm.CollsTX = uint(collsTX)
		gnm.CompressedRX = uint(compRX)
		gnm.CompressedTX = uint(compTX)
		gnm.Compressed = gnm.CompressedRX + gnm.CompressedTX
		gnm.MulticastRX = uint(multicastRX)
		gnm.CarrierTX = uint(carrierTX)

		metrics = append(metrics, gnm)
	}
}
