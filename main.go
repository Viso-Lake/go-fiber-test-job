package main

import (
	// _ "go-fiber-test-job/docs"
	"go-fiber-test-job/src/config"
	"go-fiber-test-job/src/database"
	"go-fiber-test-job/src/logger"
	"go-fiber-test-job/src/routes"
)

func init() {
	logger.InitializeLogger()
}

// @title Server API
// @version 1.0
// @description Server API

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X_API_KEY
func main() {
	config.LoadConfig()
	if config.AppConfig.IsDebug {
		logger.SetDebugLevel()
	}
	if err := database.Connect(); err != nil {
		logger.Logger.Fatal().Msg("Connect to database error. Error - " + err.Error())
	}
	app, listenAddress := routes.New()
	if err := app.Listen(listenAddress); err != nil {
		logger.Logger.Fatal().Msg("Startup error. Error - " + err.Error())
	}
}
