package testRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go-fiber-test-job/src/config"
	myLogger "go-fiber-test-job/src/logger"
	middleware "go-fiber-test-job/src/middlewares"
	accountModule "go-fiber-test-job/src/modules/account"
	cronModule "go-fiber-test-job/src/modules/cron"
)

func New() *fiber.App {
	// create app
	app := fiber.New(fiber.Config{
		AppName: config.AppConfig.AppName,
	})
	app.Use(cors.New())
	app.Use(myLogger.LogMiddleware())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.ErrorHandler())

	app.Get("/api/*", swagger.HandlerDefault)

	accountMethods := app.Group("/account")
	accountMethods.Get("/", middleware.AdminApiKeyGuard, accountModule.GetAccounts)
	accountMethods.Post("/", middleware.AdminApiKeyGuard, accountModule.CreateAccount)

	cronMethods := app.Group("/cron")
	cronMethods.Post("/account-balance", middleware.CronApiKeyGuard, cronModule.UpdateAccountsBalances)

	return app
}
