package handler

import (
	"ai-calls/internal/logger"
	"ai-calls/internal/server/app/calls"
	"ai-calls/internal/server/app/conversations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OutboundCall struct {
	From         string `json:"from"`
	To           string `json:"to" binding:"required"`
	Prompt       string `json:"prompt"`
	FirstMessage string `json:"first_message"`
}

func (h *Handler) OutboundCall(c *gin.Context) {
	var callParams OutboundCall
	if err := c.ShouldBindJSON(&callParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required fields: " + err.Error(),
		})
		return
	}

	conversation := conversations.New(h.App)
	conversationHash := conversation.GenerateHash()
	callbackTwilioUrl := h.URL.SetParam("conversation_hash", conversationHash).GetRouteUrl("https", "calls.outbound.call.twiml")
	call := calls.New(h.App)
	status, sid, err := call.CreateCall(callParams.To, callParams.From, callbackTwilioUrl)
	if err != nil {
		logger.Error("twilio error", logger.Any("api twilio", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "twilio error",
		})
		return
	}
	conversation.CreateCache(conversationHash, callParams.FirstMessage, callParams.Prompt)
	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"callSid":           sid,
		"call_status":       status,
		"conversation_hash": conversationHash,
	})
}
