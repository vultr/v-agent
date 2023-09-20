// Package api provides scrape capabilities for metrics
package api

import (
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// ExportAPI configuration for http probes
type ExportAPI struct {
	Listen string
	Port   uint

	app *fiber.App
}

// NewExporterAPI starts the fiber server
func NewExporterAPI(listen string, port uint) (*ExportAPI, error) {
	var p ExportAPI

	// Initialize engine
	app := fiber.New(fiber.Config{
		DisableKeepalive:      true,
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
	})

	prometheus := fiberprometheus.New("v-agent")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	p.app = app
	p.Listen = listen
	p.Port = port

	return &p, nil
}

// Start starts the server
func (p *ExportAPI) Start() error {
	return p.app.Listen(fmt.Sprintf("%s:%d", p.Listen, p.Port))
}

// Shutdown shuts down the server
func (p *ExportAPI) Shutdown() error {
	return p.app.Shutdown()
}
