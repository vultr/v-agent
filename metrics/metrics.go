// Package metrics provides prometheus metrics
package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	hostLoadavgLoad1        *prometheus.GaugeVec
	hostLoadavgLoad5        *prometheus.GaugeVec
	hostLoadavgLoad15       *prometheus.GaugeVec
	hostLoadavgTasksRunning *prometheus.GaugeVec
	hostLoadavgTasksTotal   *prometheus.GaugeVec
)

// NewMetrics initializes metrics
func NewMetrics() {
	// host loadavg metrics
	hostLoadavgLoad1 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load1",
			Help: "loadavg over 1 minute",
		},
		[]string{
			"hostname",
		},
	)

	hostLoadavgLoad5 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load5",
			Help: "loadavg over 5 minutes",
		},
		[]string{
			"hostname",
		},
	)

	hostLoadavgLoad15 = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_loadavg_load15",
			Help: "loadavg over 15 minutes",
		},
		[]string{
			"hostname",
		},
	)

	hostLoadavgTasksRunning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_running",
			Help: "current running tasks/processes",
		},
		[]string{
			"hostname",
		},
	)

	hostLoadavgTasksTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "v_tasks_total",
			Help: "current total tasks/processes",
		},
		[]string{
			"hostname",
		},
	)

}

// Gather gathers updates metrics
func Gather() error {
	if err := gatherHostLoadavgMetrics(); err != nil {
		return err
	}

	return nil
}

// Reset some metrics may need to be reset
func Reset() {}

func gatherHostLoadavgMetrics() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	loadavg, err := getLoadavg()
	if err != nil {
		return err
	}

	hostLoadavgLoad1.WithLabelValues(hostname).Set(loadavg.Load1)
	hostLoadavgLoad5.WithLabelValues(hostname).Set(loadavg.Load5)
	hostLoadavgLoad15.WithLabelValues(hostname).Set(loadavg.Load15)
	hostLoadavgTasksRunning.WithLabelValues(hostname).Set(float64(loadavg.TasksRunning))
	hostLoadavgTasksTotal.WithLabelValues(hostname).Set(float64(loadavg.TasksTotal))

	return nil
}
