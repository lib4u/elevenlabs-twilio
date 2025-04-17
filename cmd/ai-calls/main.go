package main

import (
	"ai-calls/internal/config"
	"ai-calls/internal/logger"
	"ai-calls/internal/server/api"
	"ai-calls/internal/server/app"
	"os"
)

func main() {
	var configPath = getConfigPath()

	config := config.Load(configPath)

	logger.Info(
		"starting application",
		logger.String("env", config.Env),
	)
	app := app.New(&config)

	api.New(app)
}

func getConfigPath() string {
	defaultConfigName := "config.yml"
	var configPath, ok = os.LookupEnv("AI_CALLS_CONF")
	if !ok {
		_, err := os.Stat(defaultConfigName)
		if os.IsNotExist(err) {
			panic("AI_CALLS_CONF Config path not found")
		}
		return defaultConfigName
	}
	return configPath
}
