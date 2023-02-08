// Package main is the main entrypoint
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vultr/v-agent/cmd/v-proxy/api"
	"github.com/vultr/v-agent/cmd/v-proxy/config"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	name string = "v-proxy"
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

	cfg, err := config.NewConfig(name)
	if err != nil {
		logger, _ := zap.NewDevelopment()
		logger.Sugar().Fatal(err)
	}

	log := zap.L().Sugar()

	e, err := api.NewHTTPAPI(cfg.Listen, cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	// http api
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("Exited http API")

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
