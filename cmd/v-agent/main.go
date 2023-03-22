// Package main is the main entrypoint
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vultr/v-agent/cmd/v-agent/api"
	"github.com/vultr/v-agent/cmd/v-agent/config"
	"github.com/vultr/v-agent/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	name    string = "v-agent"
	version string = "v0.0.13"
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

	log.Infof("version: %s", version)

	e, err := api.NewExporterAPI(cfg.Listen, cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	// /metrics endpoint
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("Exited /metrics endpoint")

				return nil
			default:
				if err := e.Start(); err != nil {
					cancel()
					log.Error(err)
				}
			}

			time.Sleep(1 * time.Second)
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

		if err := e.Shutdown(); err != nil {
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
