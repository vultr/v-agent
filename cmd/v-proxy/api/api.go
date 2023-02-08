// Package api provides scrape capabilities for metrics
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

// HTTPAPI configuration for http probes
type HTTPAPI struct {
	Listen string
	Port   uint

	app *fiber.App
}

// NewHTTPAPI starts the fiber server
func NewHTTPAPI(listen string, port uint) (*HTTPAPI, error) {
	var p HTTPAPI

	// Initialize engine
	app := fiber.New()

	stdOutLogger, err := zap.NewStdLogAt(zap.L(), zap.InfoLevel)
	if err != nil {
		return nil, err
	}

	// log configuration
	app.Use(logger.New(logger.Config{
		TimeZone:   "UTC",
		TimeFormat: "02/Jan/2006:15:04:05 -0700",
		Format:     "${ip} - - [${time}] \"${method} ${path}${queryParams} HTTP/1.1\" ${status} ${bytesSent} \"${referer}\" \"${ua}\"\n",
		Output:     stdOutLogger.Writer(),
	}))

	// rest api
	v1 := app.Group("/api/v1")
	v1.Post("/rw", V1PostRemoteWrite)

	p.app = app
	p.Listen = listen
	p.Port = port

	return &p, nil
}

// Start starts the server
func (p *HTTPAPI) Start() error {
	return p.app.Listen(fmt.Sprintf("%s:%d", p.Listen, p.Port))
}

// Shutdown shuts down the server
func (p *HTTPAPI) Shutdown() error {
	return p.app.Shutdown()
}
