package routes

import (
	"ai-calls/internal/logger"
	"ai-calls/internal/server/app"
	"ai-calls/internal/server/utils/url"
	"net/http"

	"ai-calls/internal/server/handler"

	"github.com/gin-gonic/gin"
)

type Router interface {
	LoadRoutes() *gin.Engine
	registerRoute(path string, method string, handlers ...gin.HandlerFunc)
}

type Routes struct {
	r       *gin.Engine
	URL     *url.URL
	App     *app.Application
	Handler *handler.Handler
}

func New(r *gin.Engine, u *url.URL, app *app.Application, h *handler.Handler) Router {
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)
	if len(app.Config.HTTPServer.TrustedProxies) != 0 {
		err := r.SetTrustedProxies(app.Config.HTTPServer.TrustedProxies)
		if err != nil {
			logger.Error("Error when set trusted proxy", logger.Any("error", err))
		}
	}
	var rt Router = &Routes{
		r:       r,
		URL:     u,
		App:     app,
		Handler: h,
	}
	return rt
}

func (s *Routes) LoadRoutes() *gin.Engine {

	s.registerRoute("/calls/outbound/call", http.MethodPost, s.Handler.OutboundCall)
	s.registerRoute("/calls/outbound/call/twiml", http.MethodPost, s.Handler.OutboundCallTwiml)
	s.registerRoute("/calls/outbound/call/stream", http.MethodGet, s.Handler.OutboundCallStream)
	s.registerRoute("/system/status", http.MethodGet, s.Handler.SystemStatus)
	return s.r
}

func (s *Routes) registerRoute(relativePath string, httpMethod string, handlers ...gin.HandlerFunc) {
	logger.Debug("Init route", logger.String("path", relativePath), logger.String("method", httpMethod))
	s.URL.AddRouteToUrl(relativePath)
	s.r.Handle(httpMethod, relativePath, handlers...)
}
