package handler

import (
	"ai-calls/internal/server/app/systemStatus"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SystemStatus(c *gin.Context) {
	sysStatus := systemStatus.New(h.App)
	countConnectedToWebsocket := sysStatus.GetWebSocketConnecionCount()
	c.JSON(http.StatusOK, gin.H{
		"success":                 true,
		"websocket_connect_count": countConnectedToWebsocket,
	})
}
