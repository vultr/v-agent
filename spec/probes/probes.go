// Package probes provides probes functionality via http
package probes

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	// STATE_SUCCESS when the service is in success state
	STATE_SUCCESS string = "SUCCESS" //nolint
	// STATE_FAILED when the service is in a failed state
	STATE_FAILED string = "FAILED" //nolint
	// STATE_UNKNOWN when the service is unknown (default) state
	STATE_UNKNOWN string = "UNKNOWN" //nolint
)

// ProbesAPI configuration for http probes
type ProbesAPI struct {
	Listen string
	Port   uint16

	state string

	app *fiber.App
}

// NewProbesAPI creates a new probes server
func NewProbesAPI(name, listen string, port uint16) (*ProbesAPI, error) {
	var p ProbesAPI

	// Initialize engine
	app := fiber.New(fiber.Config{
		AppName:               name,
		EnablePrintRoutes:     false,
		Prefork:               false,
		Concurrency:           50,
		ServerHeader:          name,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           30 * time.Second,
		DisableKeepalive:      true,
		DisableStartupMessage: true,
	})

	// health/ready checks
	app.Get("/healthz", p.GetHealthz)
	app.Get("/ready", p.GetReady)
	app.Get("/heap", GetHeap)

	prometheus := fiberprometheus.New(name)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	p.app = app
	p.Listen = listen
	p.Port = port
	p.state = STATE_UNKNOWN

	return &p, nil
}

// Start starts the server
func (p *ProbesAPI) Start() error {
	return p.app.Listen(fmt.Sprintf("%s:%d", p.Listen, p.Port))
}

// Shutdown shuts down the server
func (p *ProbesAPI) Shutdown() error {
	return p.app.Shutdown()
}

// Success sets state to success
func (p *ProbesAPI) Success() {
	p.state = STATE_SUCCESS
}

// Failed sets state to failed
func (p *ProbesAPI) Failed() {
	p.state = STATE_FAILED
}

// GetHealthz responds to /healthz
func (p *ProbesAPI) GetHealthz(c *fiber.Ctx) error {
	switch p.state {
	case STATE_UNKNOWN:
		c.Status(fiber.StatusContinue)

		return c.SendString(STATE_UNKNOWN)
	case STATE_SUCCESS:
		c.Status(fiber.StatusOK)

		return c.SendString(STATE_SUCCESS)
	case STATE_FAILED:
		c.Status(fiber.StatusBadRequest)

		return c.SendString(STATE_FAILED)
	}

	c.Status(fiber.StatusContinue)

	return c.SendString(STATE_UNKNOWN)
}

// GetReady responds to /ready
func (p *ProbesAPI) GetReady(c *fiber.Ctx) error {
	switch p.state {
	case STATE_UNKNOWN:
		c.Status(fiber.StatusContinue)

		return c.SendString(STATE_UNKNOWN)
	case STATE_SUCCESS:
		c.Status(fiber.StatusOK)

		return c.SendString(STATE_SUCCESS)
	case STATE_FAILED:
		c.Status(fiber.StatusBadRequest)

		return c.SendString(STATE_FAILED)
	}

	c.Status(fiber.StatusContinue)

	return c.SendString(STATE_UNKNOWN)
}

// GetHeap responds to /heap dump requests
func GetHeap(c *fiber.Ctx) error {
	log := zap.L().Sugar()

	tmpfile, err := os.CreateTemp(".", "heap-profile-*.pprof")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	log.With(
		"context", "rest",
		"source", c.IP(),
		"method", c.Method(),
		"url", c.OriginalURL(),
	).Infof("Writing heap profile to: %q", tmpfile.Name())
	if err := pprof.WriteHeapProfile(tmpfile); err != nil {
		return err
	}

	c.Status(fiber.StatusOK)

	return c.JSON(fiber.Map{
		"message": "done",
	})
}
