package middleware

import (
	"github.com/gofiber/fiber/v2"
	errorHelper "go-fiber-test-job/src/common/error-helpers"
	"go-fiber-test-job/src/config"
)

func AdminApiKeyGuard(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" || config.AppConfig.AdminXApiKey == "" || apiKey != config.AppConfig.AdminXApiKey {
		return errorHelper.RespondUnauthorizedError()
	}
	return c.Next()
}

func CronApiKeyGuard(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" || config.AppConfig.CronXApiKey == "" || apiKey != config.AppConfig.CronXApiKey {
		return errorHelper.RespondUnauthorizedError()
	}
	return c.Next()
}
