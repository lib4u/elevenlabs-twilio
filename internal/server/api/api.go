package api

import (
	"ai-calls/internal/logger"
	"ai-calls/internal/server/api/routes"
	"ai-calls/internal/server/app"
	"ai-calls/internal/server/handler"
	"ai-calls/internal/server/utils/url"
	"fmt"

	"github.com/gin-gonic/gin"
)

type API struct {
	App     *app.Application
	Handler *handler.Handler
	Router  *gin.Engine
	URL     *url.URL
}

func New(app *app.Application) {
	url := url.New(app)
	api := &API{
		App:     app,
		Handler: handler.New(app, url),
		Router:  gin.New(),
		URL:     url,
	}

	logger.Info("starting http server")
	r := routes.New(api.Router, api.URL, api.App, api.Handler)
	r.LoadRoutes()
	err := api.Router.Run(fmt.Sprintf(":%d", api.App.Config.HTTPServer.Port))
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
