package app

import (
	"ai-calls/internal/config"
	"ai-calls/internal/storage/cache"
)

type Application struct {
	Config *config.Config
	Cache  *cache.Cache
}

func New(config *config.Config) *Application {
	return &Application{
		Config: config,
		Cache:  cache.New(config),
	}
}
