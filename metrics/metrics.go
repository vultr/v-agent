// Package metrics provides prometheus metrics
package metrics

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	"github.com/vultr/v-agent/cmd/v-agent/config"
	prompb "go.buf.build/grpc/go/prometheus/prometheus"
	"go.uber.org/zap"
)

var (
	// load avg metrics
	loadavgLoad1        *prometheus.GaugeVec
	loadavgLoad5        *prometheus.GaugeVec
	loadavgLoad15       *prometheus.GaugeVec
	loadavgTasksRunning *prometheus.GaugeVec
	loadavgTasksTotal   *prometheus.GaugeVec

	// cpu metrics
	cpuCores        *prometheus.GaugeVec
	cpuUtilPct      *prometheus.GaugeVec
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
)

// NewMetrics initializes metrics
func NewMetrics() {
	// load avg metrics
	loadavgLoad1 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load1",
			Help: "loadavg over 1 minute",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	loadavgLoad5 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load5",
			Help: "loadavg over 5 minutes",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	loadavgLoad15 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load15",
			Help: "loadavg over 15 minutes",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	loadavgTasksRunning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_running",
			Help: "current running tasks/processes",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	loadavgTasksTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_total",
			Help: "current total tasks/processes",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	// cpu metrics
	cpuCores = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_cores",
			Help: "total cpu cores",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuUtilPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_util_pct",
			Help: "utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuUserPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_user_pct",
			Help: "user utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuSystemPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_system_pct",
			Help: "system utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuIOWaitPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_iowait_pct",
			Help: "iowait utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuIRQPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_irq_pct",
			Help: "irq utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuSoftIRQPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_soft_irq_pct",
			Help: "soft irq utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuStealPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_steal_pct",
			Help: "steal utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuGuestPct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_guest_pct",
			Help: "guest utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	cpuGuestNicePct = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_cpu_guest_nice_pct",
			Help: "guest nice utilization cpu percent",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	// memory metrics
	memoryTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_total",
			Help: "total memory",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	memoryFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_free",
			Help: "hosts free (non cached/buffered) memory",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	memoryCached = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_cached",
			Help: "cached memory",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	memoryBuffered = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_buffered",
			Help: "buffered memory",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	memorySwapTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_swap_total",
			Help: "total swap",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	memorySwapFree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_memory_swap_free",
			Help: "free swap",
		},
		[]string{
			"product",
			"hostname",
			"subid",
		},
	)

	// nic metrics
	nicBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes",
			Help: "nic stats: total bytes",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicBytesTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes_tx",
			Help: "nic stats: total bytes tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicBytesRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_bytes_rx",
			Help: "nic stats: total bytes rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicPackets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets",
			Help: "nic stats: total packets",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicPacketsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets_tx",
			Help: "nic stats: total packets tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicPacketsRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_packets_rx",
			Help: "nic stats: total packets rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors",
			Help: "nic stats: total errors",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicErrorsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors_tx",
			Help: "nic stats: total errors tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicErrorsRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_errors_rx",
			Help: "nic stats: total errors rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicDrop = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop",
			Help: "nic stats: total drop",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicDropTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop_tx",
			Help: "nic stats: total drop tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicDropRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_drop_rx",
			Help: "nic stats: total drop rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicFIFO = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo",
			Help: "nic stats: total fifo",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicFIFOTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo_tx",
			Help: "nic stats: total fifo tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicFIFORX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_fifo_rx",
			Help: "nic stats: total fifo rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicFrameRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_frame_rx",
			Help: "nic stats: total frame rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicCollsTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_colls_tx",
			Help: "nic stats: total colls tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicCompressed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed",
			Help: "nic stats: total compressed",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicCompressedTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed_tx",
			Help: "nic stats: total compressed tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicCompressedRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_compressed_rx",
			Help: "nic stats: total compressed rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicCarrierTX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_carrier_tx",
			Help: "nic stats: total carrier tx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

	nicMulticastRX = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_nic_multicast_rx",
			Help: "nic stats: total multicast rx",
		},
		[]string{
			"product",
			"hostname",
			"subid",
			"nic",
		},
	)

}

// Gather gathers updates metrics
func Gather() error {
	log := zap.L().Sugar()

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

	return nil
}

// RemoteWrite performs a remote write to the specified endpoint
func RemoteWrite(endpoint string) error {
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
				metricValue = in[i].Metric[0].Counter.GetValue()
			case dto.MetricType_GAUGE:
				metricValue = in[i].Metric[0].Gauge.GetValue()
			case dto.MetricType_HISTOGRAM:
				metricValue = in[i].Metric[0].Histogram.GetSampleSum()
			case dto.MetricType_SUMMARY:
				metricValue = in[i].Metric[0].Summary.GetSampleSum()
			case dto.MetricType_UNTYPED:
				metricValue = in[i].Metric[0].Untyped.GetValue()
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

// Reset some metrics may need to be reset
func Reset() {}

func gatherLoadavgMetrics() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	subid, err := config.GetSubID()
	if err != nil {
		return err
	}

	product, err := config.GetProduct()
	if err != nil {
		return err
	}

	loadavg, err := getLoadavg()
	if err != nil {
		return err
	}

	loadavgLoad1.WithLabelValues(*product, hostname, *subid).Set(loadavg.Load1)
	loadavgLoad5.WithLabelValues(*product, hostname, *subid).Set(loadavg.Load5)
	loadavgLoad15.WithLabelValues(*product, hostname, *subid).Set(loadavg.Load15)
	loadavgTasksRunning.WithLabelValues(*product, hostname, *subid).Set(float64(loadavg.TasksRunning))
	loadavgTasksTotal.WithLabelValues(*product, hostname, *subid).Set(float64(loadavg.TasksTotal))

	return nil
}

func gatherCPUMetrics() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	subid, err := config.GetSubID()
	if err != nil {
		return err
	}

	product, err := config.GetProduct()
	if err != nil {
		return err
	}

	cpuUtil, err := getCPUUtil()
	if err != nil {
		return err
	}

	cpuInUse := float64(cpuUtil.User + cpuUtil.Nice + cpuUtil.System + cpuUtil.IOWait + cpuUtil.IRQ + cpuUtil.SoftIRQ + cpuUtil.Steal + cpuUtil.Guest + cpuUtil.GuestNice)

	cpuCores.WithLabelValues(*product, hostname, *subid).Set(float64(getHostCPUs()))
	cpuUtilPct.WithLabelValues(*product, hostname, *subid).Set((cpuInUse / float64(cpuUtil.Idle)) * float64(100))
	cpuUserPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.User) / float64(cpuUtil.Idle)) * float64(100))
	cpuSystemPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.System) / float64(cpuUtil.Idle)) * float64(100))
	cpuIOWaitPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.IOWait) / float64(cpuUtil.Idle)) * float64(100))
	cpuIRQPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.IRQ) / float64(cpuUtil.Idle)) * float64(100))
	cpuSoftIRQPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.SoftIRQ) / float64(cpuUtil.Idle)) * float64(100))
	cpuStealPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.Steal) / float64(cpuUtil.Idle)) * float64(100))
	cpuGuestPct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.Guest) / float64(cpuUtil.Idle)) * float64(100))
	cpuGuestNicePct.WithLabelValues(*product, hostname, *subid).Set((float64(cpuUtil.GuestNice) / float64(cpuUtil.Idle)) * float64(100))

	return nil
}

func gatherMemoryMetrics() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	subid, err := config.GetSubID()
	if err != nil {
		return err
	}

	product, err := config.GetProduct()
	if err != nil {
		return err
	}

	memory, err := getMeminfo()
	if err != nil {
		return err
	}

	memoryTotal.WithLabelValues(*product, hostname, *subid).Set(float64(memory.MemTotal))
	memoryFree.WithLabelValues(*product, hostname, *subid).Set(float64(memory.MemFree))
	memoryCached.WithLabelValues(*product, hostname, *subid).Set(float64(memory.Cached))
	memoryBuffered.WithLabelValues(*product, hostname, *subid).Set(float64(memory.Buffers))
	memorySwapTotal.WithLabelValues(*product, hostname, *subid).Set(float64(memory.SwapTotal))
	memorySwapFree.WithLabelValues(*product, hostname, *subid).Set(float64(memory.SwapFree))

	return nil
}

func gatherNICMetrics() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	subid, err := config.GetSubID()
	if err != nil {
		return err
	}

	product, err := config.GetProduct()
	if err != nil {
		return err
	}

	nicStats, err := getNICStats()
	if err != nil {
		return err
	}

	for i := range nicStats {
		nicBytes.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].Bytes))
		nicBytesTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].BytesTX))
		nicBytesRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].BytesRX))
		nicPackets.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].Packets))
		nicPacketsTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].PacketsTX))
		nicPacketsRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].PacketsRX))
		nicErrors.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].Errors))
		nicErrorsTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].ErrorsTX))
		nicErrorsRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].ErrorsRX))
		nicDrop.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].Drop))
		nicDropTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].DropTX))
		nicDropRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].DropRX))
		nicFIFO.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].FIFO))
		nicFIFOTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].FIFOTX))
		nicFIFORX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].FIFORX))
		nicFrameRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].FrameRX))
		nicCollsTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].CollsTX))
		nicCompressed.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].Compressed))
		nicCompressedTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].CompressedTX))
		nicCompressedRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].CompressedRX))
		nicCarrierTX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].CarrierTX))
		nicMulticastRX.WithLabelValues(*product, hostname, *subid, nicStats[i].Interface).Set(float64(nicStats[i].MulticastRX))
	}

	return nil
}
