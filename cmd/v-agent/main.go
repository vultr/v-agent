// Package main is the main entrypoint
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/config"
	"github.com/vultr/v-agent/spec/metrics"
	"github.com/vultr/v-agent/spec/probes"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	name string = "v-agent"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Used for clean shutdown, catches signals
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cfg, err := config.NewConfig(name, version)
	if err != nil {
		logger, _ := zap.NewDevelopment()
		logger.Sugar().Fatal(err)
	}

	log := zap.L().Sugar()

	log.With(
		"context", name,
		"version", version,
		"commit", commit,
		"date", date,
	).Info()

	// output labels
	labels := config.GetLabels()
	for k, v := range labels {
		log.Infof("Label: %s = %s", k, v)
	}

	g, gCtx := errgroup.WithContext(ctx)

	log.With(
		"context", "v-inf",
	).Info("initializing probes api")

	probe, err := probes.NewProbesAPI(name, config.GetProbesAPIListen(), config.GetProbesAPIPort())
	if err != nil {
		log.Fatal(err)
	}

	// run http probes api
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.With(
					"context", "v-inf",
				).Info("probes: exited")

				return nil
			default:
				log.With(
					"context", "v-inf",
				).Info("probes: starting")

				if err2 := probe.Start(); err2 != nil {
					return err2
				}
			}
		}
	})

	// run metrics gatherer
	g.Go(func() error {
		runInterval := cfg.Interval
		counter := uint(0)

		metrics.NewMetrics()

		for {
			counter++

			select {
			case <-gCtx.Done():
				log.Info("Exited metrics gatherer")

				return nil
			default:
				// only run on runInterval
				if (counter % runInterval) == 0 {
					log.Infof("metrics worker: Gathering metrics")

					start := time.Now()

					if err := metrics.Gather(); err != nil {
						log.Error(err)
					}

					log.Infof("metrics worker: Gathered metrics in %s", time.Since(start).Round(time.Millisecond))

					log.Infof("metrics worker: Pushing metrics to %s", cfg.Endpoint)

					start = time.Now()

					mf, err := prometheus.DefaultGatherer.Gather()
					if err != nil {
						log.Error(err)
						continue
					}

					mf2, err := metrics.AddLabels(mf)
					if err != nil {
						log.Error(err)
						continue
					}

					tsList := metrics.GetMetricsAsTimeSeries(mf2)

					var ba *metrics.BasicAuth
					if cfg.BasicAuthUser != "" && cfg.BasicAuthPass != "" {
						ba = &metrics.BasicAuth{
							Username: cfg.BasicAuthUser,
							Password: cfg.BasicAuthPass,
						}
					}

					wc, err := metrics.NewWriteClient(cfg.Endpoint, &metrics.HTTPConfig{
						Timeout:   5 * time.Second,
						BasicAuth: ba,
					})
					if err != nil {
						log.Error(err)
						continue
					}

					if err := wc.Store(context.Background(), tsList); err != nil {
						log.Error(err)
						continue
					}

					log.Infof("metrics worker: Pushed metrics in %s", time.Since(start).Round(time.Millisecond))
				}
			}

			time.Sleep(1 * time.Second)
		}
	})

	// used for clean shutdown
	var start time.Time
	var delta time.Duration
	g.Go(func() error {
		<-gCtx.Done()
		start = time.Now()

		if err := probe.Shutdown(); err != nil {
			log.Error(err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Errorf("exit reason: %s", err)
	}

	delta = time.Since(start)

	log.Infof("Shutdown in %s", delta.Round(time.Millisecond))
}
