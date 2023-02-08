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
			"hostname",
			"subid",
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

	loadavg, err := getLoadavg()
	if err != nil {
		return err
	}

	loadavgLoad1.WithLabelValues(hostname, *subid).Set(loadavg.Load1)
	loadavgLoad5.WithLabelValues(hostname, *subid).Set(loadavg.Load5)
	loadavgLoad15.WithLabelValues(hostname, *subid).Set(loadavg.Load15)
	loadavgTasksRunning.WithLabelValues(hostname, *subid).Set(float64(loadavg.TasksRunning))
	loadavgTasksTotal.WithLabelValues(hostname, *subid).Set(float64(loadavg.TasksTotal))

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

	cpuUtil, err := getCPUUtil()
	if err != nil {
		return err
	}

	cpuInUse := float64(cpuUtil.User + cpuUtil.Nice + cpuUtil.System + cpuUtil.IOWait + cpuUtil.IRQ + cpuUtil.SoftIRQ + cpuUtil.Steal + cpuUtil.Guest + cpuUtil.GuestNice)

	cpuCores.WithLabelValues(hostname, *subid).Set(float64(getHostCPUs()))
	cpuUtilPct.WithLabelValues(hostname, *subid).Set((cpuInUse / float64(cpuUtil.Idle)) * float64(100))
	cpuUserPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.User) / float64(cpuUtil.Idle)) * float64(100))
	cpuSystemPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.System) / float64(cpuUtil.Idle)) * float64(100))
	cpuIOWaitPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.IOWait) / float64(cpuUtil.Idle)) * float64(100))
	cpuIRQPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.IRQ) / float64(cpuUtil.Idle)) * float64(100))
	cpuSoftIRQPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.SoftIRQ) / float64(cpuUtil.Idle)) * float64(100))
	cpuStealPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.Steal) / float64(cpuUtil.Idle)) * float64(100))
	cpuGuestPct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.Guest) / float64(cpuUtil.Idle)) * float64(100))
	cpuGuestNicePct.WithLabelValues(hostname, *subid).Set((float64(cpuUtil.GuestNice) / float64(cpuUtil.Idle)) * float64(100))

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

	memory, err := getMeminfo()
	if err != nil {
		return err
	}

	memoryTotal.WithLabelValues(hostname, *subid).Set(float64(memory.MemTotal))
	memoryFree.WithLabelValues(hostname, *subid).Set(float64(memory.MemFree))
	memoryCached.WithLabelValues(hostname, *subid).Set(float64(memory.Cached))
	memoryBuffered.WithLabelValues(hostname, *subid).Set(float64(memory.Buffers))
	memorySwapTotal.WithLabelValues(hostname, *subid).Set(float64(memory.SwapTotal))
	memorySwapFree.WithLabelValues(hostname, *subid).Set(float64(memory.SwapFree))

	return nil
}
