// Package api provides scrape capabilities for metrics
package api

import (
	"bytes"
	"fmt"

	"github.com/vultr/v-agent/cmd/v-proxy/config"

	"github.com/gofiber/fiber/v2"
	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

// V1PostRemoteWrite proxies remote write to mimir
func V1PostRemoteWrite(c *fiber.Ctx) error {
	log := zap.L().Sugar()

	headers := c.GetReqHeaders()
	body := c.Body()

	var buf bytes.Buffer
	if _, err := buf.Write(body); err != nil {
		log.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// verify headers, if good send to mimir
	resp, err := grequests.Post(config.GetMimirEndpoint(), &grequests.RequestOptions{
		Headers:     headers,
		RequestBody: &buf,
	})

	if err != nil {
		log.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	fmt.Println(resp.StatusCode)

	return c.SendStatus(fiber.StatusOK)
}
