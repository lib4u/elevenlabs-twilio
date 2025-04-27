package handler

import (
	"ai-calls/internal/server/app"
	"ai-calls/internal/server/utils/url"
)

type Handler struct {
	App *app.Application
	URL *url.URL
}

func New(app *app.Application, url *url.URL) *Handler {
	return &Handler{
		App: app,
		URL: url,
	}
}
