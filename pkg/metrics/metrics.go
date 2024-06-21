// Package metrics metrics collection
package metrics

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"time"

	prompb "buf.build/gen/go/prometheus/prometheus/protocolbuffers/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/vultr/v-agent/cmd/v-agent/config"
	"go.uber.org/zap"
)

var (
	vAgentVersion *prometheus.GaugeVec

	// load avg metrics
	loadavgLoad1        *prometheus.GaugeVec
	loadavgLoad5        *prometheus.GaugeVec
	loadavgLoad15       *prometheus.GaugeVec
	loadavgTasksRunning *prometheus.GaugeVec
	loadavgTasksTotal   *prometheus.GaugeVec

	// cpu metrics
	cpuCores        *prometheus.GaugeVec
	cpuUtilPct      *prometheus.GaugeVec
	cpuIdlePct      *prometheus.GaugeVec
	cpuUserPct      *prometheus.GaugeVec
	cpuSystemPct    *prometheus.GaugeVec
	cpuIOWaitPct    *prometheus.GaugeVec
	cpuIRQPct       *prometheus.GaugeVec
	cpuSoftIRQPct   *prometheus.GaugeVec
	cpuStealPct     *prometheus.GaugeVec
	cpuGuestPct     *prometheus.GaugeVec
	cpuGuestNicePct *prometheus.GaugeVec

	// memory metrics
	memoryTotal     *prometheus.GaugeVec
	memoryFree      *prometheus.GaugeVec
	memoryCached    *prometheus.GaugeVec
	memoryBuffered  *prometheus.GaugeVec
	memorySwapTotal *prometheus.GaugeVec
	memorySwapFree  *prometheus.GaugeVec

	// nic metrics
	nicBytes        *prometheus.GaugeVec
	nicBytesTX      *prometheus.GaugeVec
	nicBytesRX      *prometheus.GaugeVec
	nicPackets      *prometheus.GaugeVec
	nicPacketsTX    *prometheus.GaugeVec
	nicPacketsRX    *prometheus.GaugeVec
	nicErrors       *prometheus.GaugeVec
	nicErrorsTX     *prometheus.GaugeVec
	nicErrorsRX     *prometheus.GaugeVec
	nicDrop         *prometheus.GaugeVec
	nicDropTX       *prometheus.GaugeVec
	nicDropRX       *prometheus.GaugeVec
	nicFIFO         *prometheus.GaugeVec
	nicFIFOTX       *prometheus.GaugeVec
	nicFIFORX       *prometheus.GaugeVec
	nicFrameRX      *prometheus.GaugeVec
	nicCollsTX      *prometheus.GaugeVec
	nicCompressed   *prometheus.GaugeVec
	nicCompressedTX *prometheus.GaugeVec
	nicCompressedRX *prometheus.GaugeVec
	nicCarrierTX    *prometheus.GaugeVec
	nicMulticastRX  *prometheus.GaugeVec

	// disk metrics
	diskStatsReads                  *prometheus.GaugeVec
	diskStatsReadsMerged            *prometheus.GaugeVec
	diskStatsSectorsRead            *prometheus.GaugeVec
	diskStatsMillisecondsReading    *prometheus.GaugeVec
	diskStatsWritesCompleted        *prometheus.GaugeVec
	diskStatsWritesMerged           *prometheus.GaugeVec
	diskStatsSectorsWritten         *prometheus.GaugeVec
	diskStatsMillisecondsWriting    *prometheus.GaugeVec
	diskStatsIOsInProgress          *prometheus.GaugeVec
	diskStatsMillisecondsInIOs      *prometheus.GaugeVec
	diskStatsWeightedIOsInMS        *prometheus.GaugeVec
	diskStatsDiscards               *prometheus.GaugeVec
	diskStatsDiscardsMerged         *prometheus.GaugeVec
	diskStatsSectorsDiscarded       *prometheus.GaugeVec
	diskStatsMillisecondsDiscarding *prometheus.GaugeVec

	// filesystem metrics
	fsInodes     *prometheus.GaugeVec
	fsInodesUsed *prometheus.GaugeVec
	fsInodesUtil *prometheus.GaugeVec
	fsBytes      *prometheus.GaugeVec
	fsBytesUsed  *prometheus.GaugeVec
	fsBytesUtil  *prometheus.GaugeVec

	// kubernetes
	kubeAPIServerHealthz *prometheus.GaugeVec

	// konnectivity
	konnectivityHealthz *prometheus.GaugeVec

	// etcd
	etcdServerHealth *prometheus.GaugeVec

	// nginx-vts
	nginxVtsHealthy *prometheus.GaugeVec

	// v-cdn-agent
	vcdnAgentHealth *prometheus.GaugeVec

	// haproxy
	haproxyHealthy *prometheus.GaugeVec

	// v-dns
	vdnsHealthy *prometheus.GaugeVec

	// smart: generic
	smartPowerCycles  *prometheus.GaugeVec
	smartPowerOnHours *prometheus.GaugeVec

	// smart: nvme
	smartNvmeCriticalWarning      *prometheus.GaugeVec
	smartNvmeTemperature          *prometheus.GaugeVec
	smartNvmeAvailSpare           *prometheus.GaugeVec
	smartNvmePercentUsed          *prometheus.GaugeVec
	smartNvmeEnduranceCritWarning *prometheus.GaugeVec
	smartNvmeDataUnitsRead        *prometheus.GaugeVec
	smartNvmeDataUnitsWritten     *prometheus.GaugeVec
	smartNvmeHostReads            *prometheus.GaugeVec
	smartNvmeHostWrites           *prometheus.GaugeVec
	smartNvmeCtrlBusyTime         *prometheus.GaugeVec
	smartNvmeUnsafeShutdowns      *prometheus.GaugeVec
	smartNvmeMediaErrors          *prometheus.GaugeVec
	smartNvmeNumErrLogEntries     *prometheus.GaugeVec
	smartNvmeWarningTempTime      *prometheus.GaugeVec
	smartNvmeCritCompTime         *prometheus.GaugeVec

	// smart: sata
	smartSataPercentLifetimeRemain *prometheus.GaugeVec
	smartSataTotalLBAsWritten      *prometheus.GaugeVec
	smartSataReallocateNANDBlkCnt  *prometheus.GaugeVec
	smartSataOfflineUncorrectable  *prometheus.GaugeVec
	smartSataRawReadErrorRate      *prometheus.GaugeVec
	smartSataReportedUncorrect     *prometheus.GaugeVec
	smartSataErrorCorrectionCount  *prometheus.GaugeVec
	smartSataWriteErrorRate        *prometheus.GaugeVec
	smartSataEraseFailCount        *prometheus.GaugeVec
	smartSataProgramFailCount      *prometheus.GaugeVec
	smartSataReallocatedEventCount *prometheus.GaugeVec
	smartSataHostProgramPageCount  *prometheus.GaugeVec
	smartSataUnusedReserveNANDBlk  *prometheus.GaugeVec
	smartSataAveBlockEraseCount    *prometheus.GaugeVec
	smartSataSATAInterfacDownshift *prometheus.GaugeVec
	smartSataUnexpectPowerLossCt   *prometheus.GaugeVec
	smartSataTemperatureCelsius    *prometheus.GaugeVec
	smartSataCurrentPendingECCCnt  *prometheus.GaugeVec
	smartSataUDMACRCErrorCount     *prometheus.GaugeVec
	smartSataFTLProgramPageCount   *prometheus.GaugeVec
	smartSataSuccessRAINRecovCnt   *prometheus.GaugeVec
	smartSataPercentLifeRemaining  *prometheus.GaugeVec
	smartSataReallocatedSectorCt   *prometheus.GaugeVec
	smartSataAvailableReservdSpace *prometheus.GaugeVec
	smartSataEraseFailCountTotal   *prometheus.GaugeVec
	smartSataPowerLossCapTest      *prometheus.GaugeVec
	smartSataTotalLBAsRead         *prometheus.GaugeVec
	smartSataReadSoftErrorRate     *prometheus.GaugeVec
	smartSataEndtoEndError         *prometheus.GaugeVec
	smartSataUncorrectableErrorCnt *prometheus.GaugeVec
	smartSataCRCErrorCount         *prometheus.GaugeVec
	smartSataThermalThrottleStatus *prometheus.GaugeVec
	smartSataUnsafeShutdownCount   *prometheus.GaugeVec
	smartSataProgramFailCountChip  *prometheus.GaugeVec
	smartSataUsedRsvdBlkCntTot     *prometheus.GaugeVec
	smartSataUnusedRsvdBlkCntTot   *prometheus.GaugeVec
	smartSataProgramFailCntTotal   *prometheus.GaugeVec
	smartSataCurrentPendingSector  *prometheus.GaugeVec
	smartSataEndofLife             *prometheus.GaugeVec
	smartSataHardwareECCRecovered  *prometheus.GaugeVec
	smartSataSoftReadErrorRate     *prometheus.GaugeVec
	smartSataMediaWearoutIndicator *prometheus.GaugeVec
	smartSataDataAddressMarkErrs   *prometheus.GaugeVec
)

// NewMetrics initializes metrics
func NewMetrics() {
	vAgentVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_agent_version",
			Help: "the version of v-agent (metadata)",
		},
		[]string{},
	)

	// load avg metrics
	loadavgLoad1 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load1",
			Help: "loadavg over 1 minute",
		},
		[]string{},
	)

	loadavgLoad5 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load5",
			Help: "loadavg over 5 minutes",
		},
		[]string{},
	)

	loadavgLoad15 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load15",
			Help: "loadavg over 15 minutes",
		},
		[]string{},
	)

	loadavgTasksRunning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_running",
			Help: "current running tasks/processes",
		},
		[]string{},
	)

	loadavgTasksTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_total",
			Help: "current total tasks/processes",
		},
		[]string{},
	)

	// cpu metrics
	cpuCores = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_cores",
			Help: "total cpu cores",
		},
		[]string{},
	)

	cpuUtilPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_util_pct",
			Help: "utilization cpu percent",
		},
		[]string{},
	)

	cpuIdlePct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_idle_pct",
			Help: "idle cpu percent",
		},
		[]string{},
	)

	cpuUserPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_user_pct",
			Help: "user utilization cpu percent",
		},
		[]string{},
	)

	cpuSystemPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_system_pct",
			Help: "system utilization cpu percent",
		},
		[]string{},
	)

	cpuIOWaitPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_iowait_pct",
			Help: "iowait utilization cpu percent",
		},
		[]string{},
	)

	cpuIRQPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_irq_pct",
			Help: "irq utilization cpu percent",
		},
		[]string{},
	)

	cpuSoftIRQPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_soft_irq_pct",
			Help: "soft irq utilization cpu percent",
		},
		[]string{},
	)

	cpuStealPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_steal_pct",
			Help: "steal utilization cpu percent",
		},
		[]string{},
	)

	cpuGuestPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_guest_pct",
			Help: "guest utilization cpu percent",
		},
		[]string{},
	)

	cpuGuestNicePct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_guest_nice_pct",
			Help: "guest nice utilization cpu percent",
		},
		[]string{},
	)

	// memory metrics
	memoryTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_total",
			Help: "total memory",
		},
		[]string{},
	)

	memoryFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_free",
			Help: "hosts free (non cached/buffered) memory",
		},
		[]string{},
	)

	memoryCached = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_cached",
			Help: "cached memory",
		},
		[]string{},
	)

	memoryBuffered = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_buffered",
			Help: "buffered memory",
		},
		[]string{},
	)

	memorySwapTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_swap_total",
			Help: "total swap",
		},
		[]string{},
	)

	memorySwapFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_swap_free",
			Help: "free swap",
		},
		[]string{},
	)

	// nic metrics
	nicBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes",
			Help: "nic stats: total bytes",
		},
		[]string{
			"nic",
		},
	)

	nicBytesTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes_tx",
			Help: "nic stats: total bytes tx",
		},
		[]string{
			"nic",
		},
	)

	nicBytesRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes_rx",
			Help: "nic stats: total bytes rx",
		},
		[]string{
			"nic",
		},
	)

	nicPackets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets",
			Help: "nic stats: total packets",
		},
		[]string{
			"nic",
		},
	)

	nicPacketsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets_tx",
			Help: "nic stats: total packets tx",
		},
		[]string{
			"nic",
		},
	)

	nicPacketsRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets_rx",
			Help: "nic stats: total packets rx",
		},
		[]string{
			"nic",
		},
	)

	nicErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors",
			Help: "nic stats: total errors",
		},
		[]string{
			"nic",
		},
	)

	nicErrorsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors_tx",
			Help: "nic stats: total errors tx",
		},
		[]string{
			"nic",
		},
	)

	nicErrorsRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors_rx",
			Help: "nic stats: total errors rx",
		},
		[]string{
			"nic",
		},
	)

	nicDrop = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop",
			Help: "nic stats: total drop",
		},
		[]string{
			"nic",
		},
	)

	nicDropTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop_tx",
			Help: "nic stats: total drop tx",
		},
		[]string{
			"nic",
		},
	)

	nicDropRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop_rx",
			Help: "nic stats: total drop rx",
		},
		[]string{
			"nic",
		},
	)

	nicFIFO = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo",
			Help: "nic stats: total fifo",
		},
		[]string{
			"nic",
		},
	)

	nicFIFOTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo_tx",
			Help: "nic stats: total fifo tx",
		},
		[]string{
			"nic",
		},
	)

	nicFIFORX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo_rx",
			Help: "nic stats: total fifo rx",
		},
		[]string{
			"nic",
		},
	)

	nicFrameRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_frame_rx",
			Help: "nic stats: total frame rx",
		},
		[]string{
			"nic",
		},
	)

	nicCollsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_colls_tx",
			Help: "nic stats: total colls tx",
		},
		[]string{
			"nic",
		},
	)

	nicCompressed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed",
			Help: "nic stats: total compressed",
		},
		[]string{
			"nic",
		},
	)

	nicCompressedTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed_tx",
			Help: "nic stats: total compressed tx",
		},
		[]string{
			"nic",
		},
	)

	nicCompressedRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed_rx",
			Help: "nic stats: total compressed rx",
		},
		[]string{
			"nic",
		},
	)

	nicCarrierTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_carrier_tx",
			Help: "nic stats: total carrier tx",
		},
		[]string{
			"nic",
		},
	)

	nicMulticastRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_multicast_rx",
			Help: "nic stats: total multicast rx",
		},
		[]string{
			"nic",
		},
	)

	// disk stats
	diskStatsReads = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_read",
			Help: "disk stats: read",
		},
		[]string{
			"device",
		},
	)
	diskStatsReadsMerged = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_reads_merged",
			Help: "disk stats: merged reads",
		},
		[]string{
			"device",
		},
	)
	diskStatsSectorsRead = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_sectors_read",
			Help: "disk stats: sectors",
		},
		[]string{
			"device",
		},
	)
	diskStatsMillisecondsReading = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_ms_reading",
			Help: "disk stats: milliseconds reading",
		},
		[]string{
			"device",
		},
	)
	diskStatsWritesCompleted = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_writes_completed",
			Help: "disk stats: writes completed",
		},
		[]string{
			"device",
		},
	)
	diskStatsWritesMerged = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_writes_merged",
			Help: "disk stats: writes merged",
		},
		[]string{
			"device",
		},
	)
	diskStatsSectorsWritten = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_sectors_written",
			Help: "disk stats: sectors written",
		},
		[]string{
			"device",
		},
	)
	diskStatsMillisecondsWriting = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_ms_writing",
			Help: "disk stats: milliseconds writing",
		},
		[]string{
			"device",
		},
	)
	diskStatsIOsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_io_ip",
			Help: "disk stats: IO in progress",
		},
		[]string{
			"device",
		},
	)
	diskStatsMillisecondsInIOs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_ms_in_io",
			Help: "disk stats: milliseconds in IO",
		},
		[]string{
			"device",
		},
	)
	diskStatsWeightedIOsInMS = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_weighted_io_in_ms",
			Help: "disk stats: weighted IOs in milliseconds",
		},
		[]string{
			"device",
		},
	)
	diskStatsDiscards = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_discards",
			Help: "disk stats: discards",
		},
		[]string{
			"device",
		},
	)
	diskStatsDiscardsMerged = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_discards_merged",
			Help: "disk stats: merged discards",
		},
		[]string{
			"device",
		},
	)
	diskStatsSectorsDiscarded = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_sectors_discarded",
			Help: "disk stats: sectors discarded",
		},
		[]string{
			"device",
		},
	)
	diskStatsMillisecondsDiscarding = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_disk_stats_ms_discarding",
			Help: "disk stats: milliseconds discarding",
		},
		[]string{
			"device",
		},
	)

	// filesystem
	fsInodes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_inodes",
			Help: "filesystem inodes total",
		},
		[]string{
			"device",
			"mount",
		},
	)
	fsInodesUsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_inodes_used",
			Help: "filesystem inodes used",
		},
		[]string{
			"device",
			"mount",
		},
	)
	fsInodesUtil = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_inodes_util",
			Help: "filesystem inodes used percentage",
		},
		[]string{
			"device",
			"mount",
		},
	)
	fsBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_bytes",
			Help: "filesystem bytes total",
		},
		[]string{
			"device",
			"mount",
		},
	)
	fsBytesUsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_bytes_used",
			Help: "filesystem bytes used",
		},
		[]string{
			"device",
			"mount",
		},
	)
	fsBytesUtil = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_fs_bytes_util",
			Help: "filesystem bytes used percentage",
		},
		[]string{
			"device",
			"mount",
		},
	)

	// kubernetes
	kubeAPIServerHealthz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_kube_apiserver_healthy",
			Help: "kube-apiserver /healthz, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// konnectivity
	konnectivityHealthz = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_konnectivity_healthy",
			Help: "konnectivity /healthz, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// etcd
	etcdServerHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_etcd_healthy",
			Help: "etcd /health, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// nginx-vts
	nginxVtsHealthy = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nginx_vts_healthy",
			Help: "nginx-vts /metrics, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// v-cdn-agent
	vcdnAgentHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cdn_agent_healthy",
			Help: "v-cdn-agent /metrics, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// haproxy
	haproxyHealthy = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_haproxy_healthy",
			Help: "haproxy /metrics, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// v-dns
	vdnsHealthy = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_dns_healthy",
			Help: "v-dns /metrics, 0 = healthy, 1 = not healthy",
		},
		[]string{},
	)

	// smart
	smartPowerCycles = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_power_cycles",
			Help: "power cycles",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartPowerOnHours = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_power_on_hours",
			Help: "power on hours",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeCriticalWarning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_critical_warning",
			Help: "nvme: critical warning",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeTemperature = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_temperate",
			Help: "nvme: temperature",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeAvailSpare = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_avail_spare",
			Help: "nvme: avail spare",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmePercentUsed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_percent_used",
			Help: "nvme: percent used",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeEnduranceCritWarning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_crit_warning",
			Help: "nvme: crit warning",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeDataUnitsRead = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_data_units_read",
			Help: "nvme: data units read",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeDataUnitsWritten = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_data_units_written",
			Help: "nvme: data units written",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeHostReads = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_host_reads",
			Help: "nvme: host reads",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeHostWrites = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_host_writes",
			Help: "nvme: host writes",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeCtrlBusyTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_ctrl_busy_time",
			Help: "nvme: controller busy time",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeUnsafeShutdowns = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_unsafe_shutdowns",
			Help: "nvme: unsafe shutdowns",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeMediaErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_media_errors",
			Help: "nvme: media errors",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeNumErrLogEntries = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_err_log_entries",
			Help: "nvme: error log entries",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeWarningTempTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_warning_temp_time",
			Help: "nvme: warning temp time",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartNvmeCritCompTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_nvme_crit_comp_time",
			Help: "nvme: crit comp time",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataPercentLifetimeRemain = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_percent_lifetime_remain",
			Help: "sata: percent lifetime remain",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataTotalLBAsWritten = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_total_lbas_written",
			Help: "sata: total LBAs written",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataReallocateNANDBlkCnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_reallocate_nand_block_count",
			Help: "sata: reallocated NAND block count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataOfflineUncorrectable = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_offline_uncorrectable",
			Help: "sata: offline uncorrectable",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataRawReadErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_raw_read_error_rate",
			Help: "sata: raw read error rate",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataReportedUncorrect = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_reported_uncorrect",
			Help: "sata: reported uncorrect",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataErrorCorrectionCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_error_correction_count",
			Help: "sata: error correction count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataWriteErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_write_error_rate",
			Help: "sata: write error rate",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataEraseFailCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_erase_fail_count",
			Help: "sata: erase fail acount",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataProgramFailCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_program_fail_count",
			Help: "sata: program fail acount",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataReallocatedEventCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_reallocated_event_count",
			Help: "sata: reallocated event count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataHostProgramPageCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_host_program_page_count",
			Help: "sata: host program page count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUnusedReserveNANDBlk = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_unused_reserve_nand_block",
			Help: "sata: unused reserve NAND block",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataAveBlockEraseCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_avg_block_erase_count",
			Help: "sata: average block erase count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataSATAInterfacDownshift = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_interface_downshift",
			Help: "sata: interface downshift",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUnexpectPowerLossCt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_unexpected_powerloss_count",
			Help: "sata: unexpected powerloss count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataTemperatureCelsius = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_temperature_celsius",
			Help: "sata: temperature in celsius",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataCurrentPendingECCCnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_current_pending_ecc_count",
			Help: "sata: current pending ECC count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUDMACRCErrorCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_udma_crc_error_count",
			Help: "sata: UDMA CRC error count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataFTLProgramPageCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_ftl_program_page_count",
			Help: "sata: FTL program page count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataSuccessRAINRecovCnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_success_rain_recovery_count",
			Help: "sata: success RAIN recovery count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataPercentLifeRemaining = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_percent_life_remaining",
			Help: "sata: percent life remaining",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataReallocatedSectorCt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_reallocated_sector_count",
			Help: "sata: reallocated sector count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataAvailableReservdSpace = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_available_reserved_space",
			Help: "sata: available reserved space",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataEraseFailCountTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_erase_fail_count_total",
			Help: "sata: erase fail count total",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataPowerLossCapTest = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_power_loss_cap_test",
			Help: "sata: power loss cap test",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataTotalLBAsRead = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_total_lbas_read",
			Help: "sata: total LBAs read",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataReadSoftErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_read_soft_error_rate",
			Help: "sata: read soft error rate",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataEndtoEndError = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_end_to_end_error",
			Help: "sata: end-to-end error",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUncorrectableErrorCnt = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_uncorrectable_error_count",
			Help: "sata: uncorrectable error count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataCRCErrorCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_crc_error_count",
			Help: "sata: CRC error count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataThermalThrottleStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_thermal_throttle_status",
			Help: "sata: thermal throttle status",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUnsafeShutdownCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_unsafe_shutdown_count",
			Help: "sata: unsafe shutdown count",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataProgramFailCountChip = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_program_fail_count_chip",
			Help: "sata: program fail count chip",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUsedRsvdBlkCntTot = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_used_reserved_block_count_total",
			Help: "sata: used reserved block count total",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataUnusedRsvdBlkCntTot = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_unused_reserved_block_count_total",
			Help: "sata: unused reserved block count total",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataProgramFailCntTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_program_fail_count_total",
			Help: "sata: program fail count total",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataCurrentPendingSector = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_current_pending_sector",
			Help: "sata: current pending sector",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataEndofLife = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_end_of_life",
			Help: "sata: end of life",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataHardwareECCRecovered = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_hardware_ecc_recovered",
			Help: "sata: hardware ECC recovered",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataSoftReadErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_soft_read_error_rate",
			Help: "sata: soft read error rate",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataMediaWearoutIndicator = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_media_wearout_indicator",
			Help: "sata: media wearout indicator",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)

	smartSataDataAddressMarkErrs = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_smart_sata_data_address_mark_errors",
			Help: "sata: data address mark errors",
		},
		[]string{
			"model",
			"serial",
			"firmware",
			"device",
		},
	)
}

// Gather gathers updates metrics
func Gather() error {
	log := zap.L().Sugar()

	if err := gatherMetadataMetrics(); err != nil {
		if !errors.Is(err, config.ErrLabelNotExist) {
			return err
		} else {
			log.Warn(err)
		}
	}

	if config.LoadAvgMetricCollectionEnabled() {
		log.Info("Gathering load_avg metrics")
		if err := gatherLoadavgMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering load_avg metrics")
	}

	if config.CPUMetricCollectionEnabled() {
		log.Info("Gathering cpu metrics")
		if err := gatherCPUMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering cpu metrics")
	}

	if config.MemoryMetricCollectionEnabled() {
		log.Info("Gathering memory metrics")
		if err := gatherMemoryMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering memory metrics")
	}

	if config.NICMetricCollectionEnabled() {
		log.Info("Gathering nic metrics")
		if err := gatherNICMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering nic metrics")
	}

	if config.DiskStatsMetricCollectionEnabled() {
		log.Info("Gathering disk stats metrics")
		if err := gatherDiskMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering disk stats metrics")
	}

	if config.FileSystemMetricCollectionEnabled() {
		log.Info("Gathering file system metrics")
		if err := gatherFilesystemMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering file system metrics")
	}

	if config.KubernetesMetricCollectionEnabled() {
		log.Info("Gathering kubernetes metrics")
		if err := gatherKubernetesMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering kubernetes metrics")
	}

	if config.KonnectivityMetricCollectionEnabled() {
		log.Info("Gathering konnectivity metrics")
		if err := gatherKonnectivityMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering konnectivity metrics")
	}

	if config.EtcdMetricCollectionEnabled() {
		log.Info("Gathering etcd metrics")
		if err := gatherEtcdMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering etcd metrics")
	}

	if config.NginxVTSMetricsCollectionEnabled() {
		log.Info("Gathering nginx-vts metrics")
		if err := gatherNginxVTSMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering nginx-vts metrics")
	}

	if config.VCDNAgentMetricsCollectionEnabled() {
		log.Info("Gathering v-cdn-agent metrics")
		if err := gatherVCDNAgentMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering v-cdn-agent metrics")
	}

	if config.HAProxyMetricCollectionEnabled() {
		log.Info("Gathering haproxy metrics")
		if err := gatherHAProxyMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering haproxy metrics")
	}

	if config.CephMetricCollectionEnabled() {
		log.Info("Gathering ceph metrics")
		if err := gatherCephMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering ceph metrics")
	}

	if config.VDNSMetricsCollectionEnabled() {
		log.Info("Gathering v-dns metrics")
		if err := gatherVDNSMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering v-dns metrics")
	}

	if config.KubernetesPodsCollectionEnabled() {
		log.Info("Gathering kubernetes_pods metrics")
		if err := gatherScrapeablePodsMetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering kubernetes_pods metrics")
	}

	if config.SMARTCollectionEnabled() {
		log.Info("Gathering SMART metrics")
		if err := gatherSMARTmetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering SMART metrics")
	}

	if config.DCGMCollectionEnabled() {
		log.Info("Gathering DCGM metrics")
		if err := gatherDCGMmetrics(); err != nil {
			return err
		}
	} else {
		log.Info("Not gathering DCGM metrics")
	}

	return nil
}

// GetMetricsAsTimeSeries returns metrics as timeseries for remote write
func GetMetricsAsTimeSeries(in []*dto.MetricFamily) []*prompb.TimeSeries {
	var tsList []*prompb.TimeSeries

	for i := range in {
		metricName := *in[i].Name

		// Extract name/value of other labels and set
		for j := range in[i].Metric {
			var labels []*prompb.Label
			var ts prompb.TimeSeries
			metricType := *in[i].Type
			var metricValue float64

			// Set __name__ label for metric name
			labels = append(labels, &prompb.Label{
				Name:  "__name__",
				Value: metricName,
			})

			// For each label pair append
			for k := range in[i].Metric[j].Label {
				labels = append(labels, &prompb.Label{
					Name:  *in[i].Metric[j].Label[k].Name,
					Value: *in[i].Metric[j].Label[k].Value,
				})
			}

			// Extract metric value
			switch metricType {
			case dto.MetricType_COUNTER:
				metricValue = in[i].Metric[j].Counter.GetValue()
			case dto.MetricType_GAUGE:
				metricValue = in[i].Metric[j].Gauge.GetValue()
			case dto.MetricType_HISTOGRAM:
				metricValue = in[i].Metric[j].Histogram.GetSampleSum()
			case dto.MetricType_SUMMARY:
				metricValue = in[i].Metric[j].Summary.GetSampleSum()
			case dto.MetricType_UNTYPED:
				metricValue = in[i].Metric[j].Untyped.GetValue()
			}

			sample := &prompb.Sample{
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Value:     metricValue,
			}

			ts.Labels = labels
			ts.Samples = []*prompb.Sample{sample}

			tsList = append(tsList, &ts)
		}
	}

	return tsList
}

// ResetSMARTMetrics resets SMART metrics incase drives were swapped
func ResetSMARTMetrics() {
	smartPowerCycles.Reset()
	smartPowerOnHours.Reset()
	smartNvmeCriticalWarning.Reset()
	smartNvmeTemperature.Reset()
	smartNvmeAvailSpare.Reset()
	smartNvmePercentUsed.Reset()
	smartNvmeEnduranceCritWarning.Reset()
	smartNvmeDataUnitsRead.Reset()
	smartNvmeDataUnitsWritten.Reset()
	smartNvmeHostReads.Reset()
	smartNvmeHostWrites.Reset()
	smartNvmeCtrlBusyTime.Reset()
	smartNvmeUnsafeShutdowns.Reset()
	smartNvmeMediaErrors.Reset()
	smartNvmeNumErrLogEntries.Reset()
	smartNvmeWarningTempTime.Reset()
	smartNvmeCritCompTime.Reset()

	smartSataPercentLifetimeRemain.Reset()
	smartSataTotalLBAsWritten.Reset()
	smartSataReallocateNANDBlkCnt.Reset()
	smartSataOfflineUncorrectable.Reset()
	smartSataRawReadErrorRate.Reset()
	smartSataReportedUncorrect.Reset()
	smartSataErrorCorrectionCount.Reset()
	smartSataWriteErrorRate.Reset()
	smartSataEraseFailCount.Reset()
	smartSataProgramFailCount.Reset()
	smartSataReallocatedEventCount.Reset()
	smartSataHostProgramPageCount.Reset()
	smartSataUnusedReserveNANDBlk.Reset()
	smartSataAveBlockEraseCount.Reset()
	smartSataSATAInterfacDownshift.Reset()
	smartSataUnexpectPowerLossCt.Reset()
	smartSataTemperatureCelsius.Reset()
	smartSataCurrentPendingECCCnt.Reset()
	smartSataUDMACRCErrorCount.Reset()
	smartSataFTLProgramPageCount.Reset()
	smartSataSuccessRAINRecovCnt.Reset()
	smartSataPercentLifeRemaining.Reset()
	smartSataReallocatedSectorCt.Reset()
	smartSataAvailableReservdSpace.Reset()
	smartSataEraseFailCountTotal.Reset()
	smartSataPowerLossCapTest.Reset()
	smartSataTotalLBAsRead.Reset()
	smartSataReadSoftErrorRate.Reset()
	smartSataEndtoEndError.Reset()
	smartSataUncorrectableErrorCnt.Reset()
	smartSataCRCErrorCount.Reset()
	smartSataThermalThrottleStatus.Reset()
	smartSataUnsafeShutdownCount.Reset()
	smartSataProgramFailCountChip.Reset()
	smartSataUsedRsvdBlkCntTot.Reset()
	smartSataUnusedRsvdBlkCntTot.Reset()
	smartSataProgramFailCntTotal.Reset()
	smartSataCurrentPendingSector.Reset()
	smartSataEndofLife.Reset()
	smartSataHardwareECCRecovered.Reset()
	smartSataSoftReadErrorRate.Reset()
	smartSataMediaWearoutIndicator.Reset()
	smartSataDataAddressMarkErrs.Reset()
}

func gatherMetadataMetrics() error {
	version := config.GetVersion()

	vAgentVersion.WithLabelValues(version).Set(0)

	return nil
}

func gatherLoadavgMetrics() error {
	loadavg, err := getLoadavg()
	if err != nil {
		return err
	}

	loadavgLoad1.WithLabelValues().Set(loadavg.Load1)
	loadavgLoad5.WithLabelValues().Set(loadavg.Load5)
	loadavgLoad15.WithLabelValues().Set(loadavg.Load15)
	loadavgTasksRunning.WithLabelValues().Set(float64(loadavg.TasksRunning))
	loadavgTasksTotal.WithLabelValues().Set(float64(loadavg.TasksTotal))

	return nil
}

func gatherCPUMetrics() error {
	cpuUtil, err := getCPUUtil()
	if err != nil {
		return err
	}

	cpuTotalTime := float64(cpuUtil.User + cpuUtil.Nice + cpuUtil.System + cpuUtil.Idle + cpuUtil.IOWait + cpuUtil.IRQ + cpuUtil.SoftIRQ + cpuUtil.Steal + cpuUtil.Guest + cpuUtil.GuestNice)

	idleTime := float64(cpuUtil.Idle) / cpuTotalTime
	inUseTime := 1 - idleTime
	userTime := float64(cpuUtil.User) / cpuTotalTime
	systemTime := float64(cpuUtil.System) / cpuTotalTime
	iowaitTime := float64(cpuUtil.IOWait) / cpuTotalTime
	irqTime := float64(cpuUtil.IRQ) / cpuTotalTime
	sirqTime := float64(cpuUtil.SoftIRQ) / cpuTotalTime
	stealTime := float64(cpuUtil.Steal) / cpuTotalTime
	guestTime := float64(cpuUtil.Guest) / cpuTotalTime
	guestNiceTime := float64(cpuUtil.GuestNice) / cpuTotalTime

	cpuCores.WithLabelValues().Set(float64(getHostCPUs()))
	cpuUtilPct.WithLabelValues().Set(inUseTime * float64(100))          //nolint
	cpuIdlePct.WithLabelValues().Set(idleTime * float64(100))           //nolint
	cpuUserPct.WithLabelValues().Set(userTime * float64(100))           //nolint
	cpuSystemPct.WithLabelValues().Set(systemTime * float64(100))       //nolint
	cpuIOWaitPct.WithLabelValues().Set(iowaitTime * float64(100))       //nolint
	cpuIRQPct.WithLabelValues().Set(irqTime * float64(100))             //nolint
	cpuSoftIRQPct.WithLabelValues().Set(sirqTime * float64(100))        //nolint
	cpuStealPct.WithLabelValues().Set(stealTime * float64(100))         //nolint
	cpuGuestPct.WithLabelValues().Set(guestTime * float64(100))         //nolint
	cpuGuestNicePct.WithLabelValues().Set(guestNiceTime * float64(100)) //nolint

	return nil
}

func gatherMemoryMetrics() error {
	memory, err := getMeminfo()
	if err != nil {
		return err
	}

	memoryTotal.WithLabelValues().Set(float64(memory.MemTotal))
	memoryFree.WithLabelValues().Set(float64(memory.MemFree))
	memoryCached.WithLabelValues().Set(float64(memory.Cached))
	memoryBuffered.WithLabelValues().Set(float64(memory.Buffers))
	memorySwapTotal.WithLabelValues().Set(float64(memory.SwapTotal))
	memorySwapFree.WithLabelValues().Set(float64(memory.SwapFree))

	return nil
}

func gatherNICMetrics() error {
	nicStats, err := getNICStats()
	if err != nil {
		return err
	}

	for i := range nicStats {
		nicBytes.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].Bytes))
		nicBytesTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].BytesTX))
		nicBytesRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].BytesRX))
		nicPackets.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].Packets))
		nicPacketsTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].PacketsTX))
		nicPacketsRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].PacketsRX))
		nicErrors.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].Errors))
		nicErrorsTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].ErrorsTX))
		nicErrorsRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].ErrorsRX))
		nicDrop.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].Drop))
		nicDropTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].DropTX))
		nicDropRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].DropRX))
		nicFIFO.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].FIFO))
		nicFIFOTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].FIFOTX))
		nicFIFORX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].FIFORX))
		nicFrameRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].FrameRX))
		nicCollsTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].CollsTX))
		nicCompressed.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].Compressed))
		nicCompressedTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].CompressedTX))
		nicCompressedRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].CompressedRX))
		nicCarrierTX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].CarrierTX))
		nicMulticastRX.WithLabelValues(nicStats[i].Interface).Set(float64(nicStats[i].MulticastRX))
	}

	return nil
}

func gatherDiskMetrics() error {
	diskStats, err := getDiskStatsUtil()
	if err != nil {
		return err
	}

	for i := range diskStats {
		diskStatsReads.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].Reads))
		diskStatsReadsMerged.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].ReadsMerged))
		diskStatsSectorsRead.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].SectorsRead))
		diskStatsMillisecondsReading.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].MillisecondsReading))
		diskStatsWritesCompleted.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].WritesCompleted))
		diskStatsWritesMerged.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].WritesMerged))
		diskStatsSectorsWritten.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].SectorsWritten))
		diskStatsMillisecondsWriting.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].MillisecondsWriting))
		diskStatsIOsInProgress.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].IOsInProgress))
		diskStatsMillisecondsInIOs.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].MillisecondsInIOs))
		diskStatsWeightedIOsInMS.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].WeightedIOsInMS))
		diskStatsDiscards.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].Discards))
		diskStatsDiscardsMerged.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].DiscardsMerged))
		diskStatsSectorsDiscarded.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].SectorsDiscarded))
		diskStatsMillisecondsDiscarding.WithLabelValues(diskStats[i].Device).Set(float64(diskStats[i].MillisecondsDiscarding))
	}

	return nil
}

func gatherFilesystemMetrics() error {
	fsStats, err := getFilesystemUtil()
	if err != nil {
		return err
	}

	for i := range fsStats {
		fsInodes.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(float64(fsStats[i].Inodes))
		fsInodesUsed.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(float64(fsStats[i].InodesUsed))
		fsInodesUtil.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(fsStats[i].InodesUtil)
		fsBytes.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(float64(fsStats[i].BytesTotal))
		fsBytesUsed.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(float64(fsStats[i].BytesUsed))
		fsBytesUtil.WithLabelValues(fsStats[i].Device, fsStats[i].Mount).Set(fsStats[i].BytesUtil)
	}

	return nil
}

func gatherKubernetesMetrics() error {
	log := zap.L().Sugar()

	if err := DoKubeAPIServerHealthCheck(); err != nil {
		log.Error(err)

		kubeAPIServerHealthz.WithLabelValues().Set(float64(1))
	} else {
		kubeAPIServerHealthz.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeKubeAPIServerMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherKonnectivityMetrics() error {
	log := zap.L().Sugar()

	if err := DoKonnectivityHealthCheck(); err != nil {
		log.Error(err)

		konnectivityHealthz.WithLabelValues().Set(float64(1))
	} else {
		konnectivityHealthz.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeKonnectivityMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherEtcdMetrics() error {
	log := zap.L().Sugar()

	if err := DoEtcdHealthCheck(); err != nil {
		log.Error(err)

		etcdServerHealth.WithLabelValues().Set(float64(1))
	} else {
		etcdServerHealth.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeEtcdMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherNginxVTSMetrics() error {
	log := zap.L().Sugar()

	if err := DoNginxVTSHealthCheck(); err != nil {
		log.Error(err)

		nginxVtsHealthy.WithLabelValues().Set(float64(1))
	} else {
		nginxVtsHealthy.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeNginxVTSMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherVCDNAgentMetrics() error {
	log := zap.L().Sugar()

	if err := DoVCDNAgentHealthCheck(); err != nil {
		log.Error(err)

		vcdnAgentHealth.WithLabelValues().Set(float64(1))
	} else {
		vcdnAgentHealth.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeVCDNAgentMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherHAProxyMetrics() error {
	log := zap.L().Sugar()

	if err := DoHAProxyHealthCheck(); err != nil {
		log.Error(err)

		haproxyHealthy.WithLabelValues().Set(float64(1))
	} else {
		haproxyHealthy.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeHAProxyMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherCephMetrics() error {
	if err := ScrapeCephMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherVDNSMetrics() error {
	log := zap.L().Sugar()

	if err := DoVDNSHealthCheck(); err != nil {
		log.Error(err)

		vdnsHealthy.WithLabelValues().Set(float64(1))
	} else {
		vdnsHealthy.WithLabelValues().Set(float64(0))
	}

	if err := ScrapeVDNSMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherScrapeablePodsMetrics() error {
	if err := ScrapeKubernetesPods(); err != nil {
		return err
	}

	return nil
}

func gatherSMARTmetrics() error {
	ResetSMARTMetrics()

	if err := ScrapeSMARTMetrics(); err != nil {
		return err
	}

	return nil
}

func gatherDCGMmetrics() error {
	if err := ScrapeDCGMMetrics(); err != nil {
		return err
	}

	return nil
}

// removeMetadata removes the "comment" lines (#) from scraped output
//
// required for some metrics (ceph) as they output duplicate HELP sections which break the parser
func removeMetadata(data []byte) []byte {
	var retBytes bytes.Buffer

	b := bytes.NewBuffer(data)

	sc := bufio.NewScanner(b)
	for sc.Scan() {
		l := sc.Text()

		if strings.HasPrefix(l, "#") {
			continue
		} else {
			retBytes.WriteString(l)
			retBytes.WriteString("\n")
		}
	}

	return retBytes.Bytes()
}

func parseMetrics(data []byte) ([]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	var mf2 []*dto.MetricFamily

	buf := bytes.NewBuffer(data)

	mf, err := parser.TextToMetricFamilies(buf)
	if err != nil {
		return nil, err
	}

	// convert to workable type []*dto.MetricFamily
	for _, v := range mf {
		mf2 = append(mf2, v)
	}

	mf2, err = AddLabels(mf2)
	if err != nil {
		return nil, err
	}

	return mf2, nil
}

// AddLabels adds labels to metrics
//
// very important for scraped metrics, otherwise they'll be sent without essential labels (subid, etc)
func AddLabels(metrics []*dto.MetricFamily) ([]*dto.MetricFamily, error) {
	arbitraryLabels := config.GetLabels()

	for _, v := range metrics {
		metricType := v.GetType()

		for _, vv := range v.Metric {
			switch metricType {
			case dto.MetricType_COUNTER:
				if vv.Counter == nil {
					continue
				}

				labels := vv.GetLabel()

				chxLabels := make(map[string]bool)

				for i := range arbitraryLabels {
					// true if label exists
					chxLabels[i] = false

					for j := range labels {
						if *labels[j].Name == i {
							chxLabels[i] = true
							break
						}
					}
				}

				// add missing labels
				for i := range chxLabels {
					if !chxLabels[i] {
						k := i
						v := arbitraryLabels[i]
						vv.Label = append(vv.Label,
							&dto.LabelPair{
								Name:  &k,
								Value: &v,
							},
						)
					}
				}
			case dto.MetricType_GAUGE:
				if vv.Gauge == nil {
					continue
				}

				labels := vv.GetLabel()

				chxLabels := make(map[string]bool)

				for i := range arbitraryLabels {
					// true if label exists
					chxLabels[i] = false

					for j := range labels {
						if *labels[j].Name == i {
							chxLabels[i] = true
							break
						}
					}
				}

				// add missing labels
				for i := range chxLabels {
					if !chxLabels[i] {
						k := i
						v := arbitraryLabels[i]
						vv.Label = append(vv.Label,
							&dto.LabelPair{
								Name:  &k,
								Value: &v,
							},
						)
					}
				}
			case dto.MetricType_UNTYPED:
				if vv.Untyped == nil {
					continue
				}

				labels := vv.GetLabel()

				chxLabels := make(map[string]bool)

				for i := range arbitraryLabels {
					// true if label exists
					chxLabels[i] = false

					for j := range labels {
						if *labels[j].Name == i {
							chxLabels[i] = true
							break
						}
					}
				}

				// add missing labels
				for i := range chxLabels {
					if !chxLabels[i] {
						k := i
						v := arbitraryLabels[i]
						vv.Label = append(vv.Label,
							&dto.LabelPair{
								Name:  &k,
								Value: &v,
							},
						)
					}
				}
			case dto.MetricType_SUMMARY:
				if vv.Summary == nil {
					continue
				}

				labels := vv.GetLabel()

				chxLabels := make(map[string]bool)

				for i := range arbitraryLabels {
					// true if label exists
					chxLabels[i] = false

					for j := range labels {
						if *labels[j].Name == i {
							chxLabels[i] = true
							break
						}
					}
				}

				// add missing labels
				for i := range chxLabels {
					if !chxLabels[i] {
						k := i
						v := arbitraryLabels[i]
						vv.Label = append(vv.Label,
							&dto.LabelPair{
								Name:  &k,
								Value: &v,
							},
						)
					}
				}
			case dto.MetricType_HISTOGRAM:
				if vv.Histogram == nil {
					continue
				}

				labels := vv.GetLabel()

				chxLabels := make(map[string]bool)

				for i := range arbitraryLabels {
					// true if label exists
					chxLabels[i] = false

					for j := range labels {
						if *labels[j].Name == i {
							chxLabels[i] = true
							break
						}
					}
				}

				// add missing labels
				for i := range chxLabels {
					if !chxLabels[i] {
						k := i
						v := arbitraryLabels[i]
						vv.Label = append(vv.Label,
							&dto.LabelPair{
								Name:  &k,
								Value: &v,
							},
						)
					}
				}
			}
		}
	}

	return metrics, nil
}
