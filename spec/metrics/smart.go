// Package metrics provides prometheus metrics
package metrics

import (
	"github.com/vultr/v-agent/cmd/v-agent/config"
	"go.uber.org/zap"

	"github.com/anatol/smart.go"
)

type SMART struct {
	Device string

	NVMe bool
	SATA bool
	SCSI bool

	Model    string
	Serial   string
	Firmware string

	// generic
	PowerCycles  uint64
	PowerOnHours uint64

	// nvme
	CriticalWarning      uint8
	Temperature          uint16
	AvailSpare           uint8
	SpareThresh          uint8
	PercentUsed          uint8
	EnduranceCritWarning uint8

	DataUnitsRead    uint64 // SHOULD BE uint128 but go doesnt support
	DataUnitsWritten uint64 // SHOULD BE uint128 but go doesnt support
	HostReads        uint64 // SHOULD BE uint128 but go doesnt support
	HostWrites       uint64 // SHOULD BE uint128 but go doesnt support
	CtrlBusyTime     uint64 // SHOULD BE uint128 but go doesnt support
	UnsafeShutdowns  uint64 // SHOULD BE uint128 but go doesnt support
	MediaErrors      uint64 // SHOULD BE uint128 but go doesnt support
	NumErrLogEntries uint64 // SHOULD BE uint128 but go doesnt support
	WarningTempTime  uint32
	CritCompTime     uint32

	// sata
	PercentLifetimeRemain uint64
	TotalLBAsWritten      uint64
	ReallocateNANDBlkCnt  uint64
	OfflineUncorrectable  uint64
	RawReadErrorRate      uint64
	ReportedUncorrect     uint64
	ErrorCorrectionCount  uint64
	WriteErrorRate        uint64
	EraseFailCount        uint64
	ProgramFailCount      uint64
	ReallocatedEventCount uint64
	HostProgramPageCount  uint64
	UnusedReserveNANDBlk  uint64
	AveBlockEraseCount    uint64
	SATAInterfacDownshift uint64
	UnexpectPowerLossCt   uint64
	TemperatureCelsius    uint64
	CurrentPendingECCCnt  uint64
	UDMACRCErrorCount     uint64
	FTLProgramPageCount   uint64
	SuccessRAINRecovCnt   uint64
	PercentLifeRemaining  uint64
	ReallocatedSectorCt   uint64
	AvailableReservdSpace uint64
	EraseFailCountTotal   uint64
	PowerLossCapTest      uint64
	TotalLBAsRead         uint64
	ReadSoftErrorRate     uint64
	EndtoEndError         uint64
	UncorrectableErrorCnt uint64
	CRCErrorCount         uint64
	ThermalThrottleStatus uint64
	UnsafeShutdownCount   uint64
	ProgramFailCountChip  uint64
	UsedRsvdBlkCntTot     uint64
	UnusedRsvdBlkCntTot   uint64
	ProgramFailCntTotal   uint64
	CurrentPendingSector  uint64
	EndofLife             uint64
	HardwareECCRecovered  uint64
	SoftReadErrorRate     uint64
	MediaWearoutIndicator uint64
	DataAddressMarkErrs   uint64
}

// ProbeSMARTBlockDevice returns a SMART struct with the device and the SMART attributes
//
// Supports NVMe and SATA
func ProbeSMARTBlockDevice(procDev string) (*SMART, error) {
	log := zap.L().Sugar()

	var smartDev SMART
	smartDev.Device = procDev

	dev, err := smart.Open(procDev)
	if err != nil {
		return nil, err
	}
	defer dev.Close()

	switch sm := dev.(type) {
	case *smart.SataDevice:
		smartDev.SATA = true

		i, err := sm.Identify()
		if err != nil {
			return nil, err
		}

		smartDev.Model = i.ModelNumber()
		smartDev.Serial = i.SerialNumber()
		smartDev.Firmware = i.FirmwareRevision()

		ata, err := sm.ReadSMARTData()
		if err != nil {
			return nil, err
		}

		for i := range ata.Attrs {
			switch ata.Attrs[i].Name {
			case "Power_Cycle_Count":
				smartDev.PowerCycles = ata.Attrs[i].ValueRaw
			case "Percent_Lifetime_Remain":
				smartDev.PercentLifetimeRemain = ata.Attrs[i].ValueRaw
			case "Total_LBAs_Written":
				smartDev.TotalLBAsWritten = ata.Attrs[i].ValueRaw
			case "Reallocate_NAND_Blk_Cnt":
				smartDev.ReallocateNANDBlkCnt = ata.Attrs[i].ValueRaw
			case "Power_On_Hours":
				smartDev.PowerOnHours = ata.Attrs[i].ValueRaw
			case "Offline_Uncorrectable":
				smartDev.OfflineUncorrectable = ata.Attrs[i].ValueRaw
			case "Raw_Read_Error_Rate":
				smartDev.RawReadErrorRate = ata.Attrs[i].ValueRaw
			case "Reported_Uncorrect":
				smartDev.ReportedUncorrect = ata.Attrs[i].ValueRaw
			case "Error_Correction_Count":
				smartDev.ErrorCorrectionCount = ata.Attrs[i].ValueRaw
			case "Write_Error_Rate":
				smartDev.WriteErrorRate = ata.Attrs[i].ValueRaw
			case "Erase_Fail_Count":
				smartDev.EraseFailCount = ata.Attrs[i].ValueRaw
			case "Program_Fail_Count":
				smartDev.ProgramFailCount = ata.Attrs[i].ValueRaw
			case "Reallocated_Event_Count":
				smartDev.ReallocatedEventCount = ata.Attrs[i].ValueRaw
			case "Host_Program_Page_Count":
				smartDev.HostProgramPageCount = ata.Attrs[i].ValueRaw
			case "Unused_Reserve_NAND_Blk":
				smartDev.UnusedReserveNANDBlk = ata.Attrs[i].ValueRaw
			case "Ave_Block-Erase_Count":
				smartDev.AveBlockEraseCount = ata.Attrs[i].ValueRaw
			case "SATA_Interfac_Downshift":
				smartDev.SATAInterfacDownshift = ata.Attrs[i].ValueRaw
			case "Unexpect_Power_Loss_Ct":
				smartDev.UnexpectPowerLossCt = ata.Attrs[i].ValueRaw
			case "Temperature_Celsius":
				smartDev.TemperatureCelsius = ata.Attrs[i].ValueRaw
			case "Current_Pending_ECC_Cnt":
				smartDev.CurrentPendingECCCnt = ata.Attrs[i].ValueRaw
			case "UDMA_CRC_Error_Count":
				smartDev.UDMACRCErrorCount = ata.Attrs[i].ValueRaw
			case "FTL_Program_Page_Count":
				smartDev.FTLProgramPageCount = ata.Attrs[i].ValueRaw
			case "Success_RAIN_Recov_Cnt":
				smartDev.SuccessRAINRecovCnt = ata.Attrs[i].ValueRaw
			case "Percent_Life_Remaining":
				smartDev.PercentLifeRemaining = ata.Attrs[i].ValueRaw
			case "Reallocated_Sector_Ct":
				smartDev.ReallocatedSectorCt = ata.Attrs[i].ValueRaw
			case "Available_Reservd_Space":
				smartDev.AvailableReservdSpace = ata.Attrs[i].ValueRaw
			case "Erase_Fail_Count_Total":
				smartDev.EraseFailCountTotal = ata.Attrs[i].ValueRaw
			case "Power_Loss_Cap_Test":
				smartDev.PowerLossCapTest = ata.Attrs[i].ValueRaw
			case "Total_LBAs_Read":
				smartDev.TotalLBAsRead = ata.Attrs[i].ValueRaw
			case "Read_Soft_Error_Rate":
				smartDev.ReadSoftErrorRate = ata.Attrs[i].ValueRaw
			case "End-to-End_Error":
				smartDev.EndtoEndError = ata.Attrs[i].ValueRaw
			case "Uncorrectable_Error_Cnt":
				smartDev.UncorrectableErrorCnt = ata.Attrs[i].ValueRaw
			case "CRC_Error_Count":
				smartDev.CRCErrorCount = ata.Attrs[i].ValueRaw
			case "Thermal_Throttle_Status":
				smartDev.ThermalThrottleStatus = ata.Attrs[i].ValueRaw
			case "Unsafe_Shutdown_Count":
				smartDev.UnsafeShutdownCount = ata.Attrs[i].ValueRaw
			case "Program_Fail_Count_Chip":
				smartDev.ProgramFailCountChip = ata.Attrs[i].ValueRaw
			case "Used_Rsvd_Blk_Cnt_Tot":
				smartDev.UsedRsvdBlkCntTot = ata.Attrs[i].ValueRaw
			case "Unused_Rsvd_Blk_Cnt_Tot":
				smartDev.UnusedRsvdBlkCntTot = ata.Attrs[i].ValueRaw
			case "Program_Fail_Cnt_Total":
				smartDev.ProgramFailCntTotal = ata.Attrs[i].ValueRaw
			case "Current_Pending_Sector":
				smartDev.CurrentPendingSector = ata.Attrs[i].ValueRaw
			case "End_of_Life":
				smartDev.EndofLife = ata.Attrs[i].ValueRaw
			case "Hardware_ECC_Recovered":
				smartDev.HardwareECCRecovered = ata.Attrs[i].ValueRaw
			case "Soft_Read_Error_Rate":
				smartDev.SoftReadErrorRate = ata.Attrs[i].ValueRaw
			case "Media_Wearout_Indicator":
				smartDev.MediaWearoutIndicator = ata.Attrs[i].ValueRaw
			case "Data_Address_Mark_Errs":
				smartDev.DataAddressMarkErrs = ata.Attrs[i].ValueRaw
			default:
				if ata.Attrs[i].Name != "" {
					log.With(
						"device", smartDev.Device,
						"type", "SATA",
					).Warnf("SMART attribute %q unknown", ata.Attrs[i].Name)
				}
			}
		}

	case *smart.ScsiDevice:
		smartDev.SCSI = true

	case *smart.NVMeDevice:
		smartDev.NVMe = true

		i, _, err := sm.Identify()
		if err != nil {
			return nil, err
		}

		smartDev.Model = i.ModelNumber()
		smartDev.Serial = i.SerialNumber()
		smartDev.Firmware = i.FirmwareRev()

		gen, err := sm.ReadGenericAttributes()
		if err != nil {
			return nil, err
		}

		smartDev.PowerCycles = gen.PowerCycles
		smartDev.PowerOnHours = gen.PowerOnHours

		nvmeLog, err := sm.ReadSMART()
		if err != nil {
			return nil, err
		}

		smartDev.CriticalWarning = nvmeLog.CritWarning
		smartDev.Temperature = nvmeLog.Temperature
		smartDev.AvailSpare = nvmeLog.AvailSpare
		smartDev.SpareThresh = nvmeLog.SpareThresh
		smartDev.PercentUsed = nvmeLog.PercentUsed
		smartDev.EnduranceCritWarning = nvmeLog.EnduranceCritWarning
		smartDev.DataUnitsRead = nvmeLog.DataUnitsRead.Val[0]
		smartDev.DataUnitsWritten = nvmeLog.DataUnitsRead.Val[0]
		smartDev.HostReads = nvmeLog.HostReads.Val[0]
		smartDev.HostWrites = nvmeLog.HostWrites.Val[0]
		smartDev.CtrlBusyTime = nvmeLog.CtrlBusyTime.Val[0]
		smartDev.UnsafeShutdowns = nvmeLog.UnsafeShutdowns.Val[0]
		smartDev.MediaErrors = nvmeLog.MediaErrors.Val[0]
		smartDev.NumErrLogEntries = nvmeLog.NumErrLogEntries.Val[0]
		smartDev.WarningTempTime = nvmeLog.WarningTempTime
		smartDev.CritCompTime = nvmeLog.CritCompTime
	}

	return &smartDev, nil
}

// ScrapeSMARTMetrics probes block device metrics
func ScrapeSMARTMetrics() error {
	log := zap.L().Sugar()

	blockDevices := config.GetSMARTBlockDevices()

	// scrape SMART data for block devices
	var smartData []*SMART
	for i := range blockDevices {
		smt, err := ProbeSMARTBlockDevice(blockDevices[i])
		if err != nil {
			log.With(
				"device", blockDevices[i],
			).Warn(err.Error())

			continue
		}

		smartData = append(smartData, smt)
	}

	if len(smartData) == 0 {
		log.Warnf("no SMART data found from block devices: %+v", blockDevices)
	}

	// set metrics
	for i := range smartData {
		smartPowerCycles.WithLabelValues(
			smartData[i].Model,
			smartData[i].Serial,
			smartData[i].Firmware,
			smartData[i].Device,
		).Set(float64(smartData[i].PowerCycles))

		smartPowerOnHours.WithLabelValues(
			smartData[i].Model,
			smartData[i].Serial,
			smartData[i].Firmware,
			smartData[i].Device,
		).Set(float64(smartData[i].PowerOnHours))

		if smartData[i].NVMe {
			smartNvmeCriticalWarning.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CriticalWarning))

			smartNvmeTemperature.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].Temperature))

			smartNvmeAvailSpare.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].AvailSpare))

			smartNvmePercentUsed.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].PercentUsed))

			smartNvmeEnduranceCritWarning.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].EnduranceCritWarning))

			smartNvmeDataUnitsRead.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].DataUnitsRead))

			smartNvmeDataUnitsWritten.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].DataUnitsWritten))

			smartNvmeHostReads.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].HostReads))

			smartNvmeHostWrites.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].HostWrites))

			smartNvmeCtrlBusyTime.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CtrlBusyTime))

			smartNvmeUnsafeShutdowns.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UnsafeShutdowns))

			smartNvmeMediaErrors.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].MediaErrors))

			smartNvmeNumErrLogEntries.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].NumErrLogEntries))

			smartNvmeWarningTempTime.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].WarningTempTime))

			smartNvmeCritCompTime.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CritCompTime))
		} else if smartData[i].SATA {
			smartSataPercentLifetimeRemain.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].PercentLifetimeRemain))

			smartSataTotalLBAsWritten.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].TotalLBAsWritten))

			smartSataReallocateNANDBlkCnt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ReallocateNANDBlkCnt))

			smartSataOfflineUncorrectable.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].OfflineUncorrectable))

			smartSataRawReadErrorRate.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].RawReadErrorRate))

			smartSataReportedUncorrect.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ReportedUncorrect))

			smartSataErrorCorrectionCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ErrorCorrectionCount))

			smartSataWriteErrorRate.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].WriteErrorRate))

			smartSataEraseFailCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].EraseFailCount))

			smartSataProgramFailCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ProgramFailCount))

			smartSataReallocatedEventCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ReallocatedEventCount))

			smartSataHostProgramPageCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].HostProgramPageCount))

			smartSataUnusedReserveNANDBlk.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UnusedReserveNANDBlk))

			smartSataAveBlockEraseCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].AveBlockEraseCount))

			smartSataSATAInterfacDownshift.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].SATAInterfacDownshift))

			smartSataUnexpectPowerLossCt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UnexpectPowerLossCt))

			smartSataTemperatureCelsius.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].TemperatureCelsius))

			smartSataCurrentPendingECCCnt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CurrentPendingECCCnt))

			smartSataUDMACRCErrorCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UDMACRCErrorCount))

			smartSataFTLProgramPageCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].FTLProgramPageCount))

			smartSataSuccessRAINRecovCnt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].SuccessRAINRecovCnt))

			smartSataPercentLifeRemaining.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].PercentLifeRemaining))

			smartSataReallocatedSectorCt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ReallocatedSectorCt))

			smartSataAvailableReservdSpace.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].AvailableReservdSpace))

			smartSataEraseFailCountTotal.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].EraseFailCountTotal))

			smartSataPowerLossCapTest.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].PowerLossCapTest))

			smartSataTotalLBAsRead.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].TotalLBAsRead))

			smartSataReadSoftErrorRate.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ReadSoftErrorRate))

			smartSataEndtoEndError.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].EndtoEndError))

			smartSataUncorrectableErrorCnt.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UncorrectableErrorCnt))

			smartSataCRCErrorCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CRCErrorCount))

			smartSataThermalThrottleStatus.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ThermalThrottleStatus))

			smartSataUnsafeShutdownCount.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UnsafeShutdownCount))

			smartSataProgramFailCountChip.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ProgramFailCountChip))

			smartSataUsedRsvdBlkCntTot.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UsedRsvdBlkCntTot))

			smartSataUnusedRsvdBlkCntTot.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].UnusedRsvdBlkCntTot))

			smartSataProgramFailCntTotal.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].ProgramFailCntTotal))

			smartSataCurrentPendingSector.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].CurrentPendingSector))

			smartSataEndofLife.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].EndofLife))

			smartSataHardwareECCRecovered.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].HardwareECCRecovered))

			smartSataSoftReadErrorRate.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].SoftReadErrorRate))

			smartSataMediaWearoutIndicator.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].MediaWearoutIndicator))

			smartSataDataAddressMarkErrs.WithLabelValues(
				smartData[i].Model,
				smartData[i].Serial,
				smartData[i].Firmware,
				smartData[i].Device,
			).Set(float64(smartData[i].DataAddressMarkErrs))
		} else if smartData[i].SCSI {
			log.With(
				"device", smartData[i].Device,
			).Warn("block device is SCSI not implemented")
		} else {
			log.With(
				"device", smartData[i].Device,
			).Warn("device is unknown: not SATA, NVMe, or SCSI")
		}
	}

	return nil
}
