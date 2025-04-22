package handler

import (
	"ai-calls/internal/logger"
	"ai-calls/internal/server/app/calls"
	"ai-calls/internal/server/app/conversations"
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OutboundCall struct {
	From             string         `json:"from"`
	DynamicVariables map[string]any `json:"dynamic_variables"`
	To               string         `json:"to" binding:"required"`
	Prompt           string         `json:"prompt"`
	FirstMessage     string         `json:"first_message"`
}

type Response struct {
	XMLName xml.Name `xml:"Response"`
	Connect Connect  `xml:"Connect"`
}

type Connect struct {
	Stream Stream `xml:"Stream"`
}

type Stream struct {
	URL        string      `xml:"url,attr"`
	Parameters []Parameter `xml:"Parameter"`
}

type Parameter struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (h *Handler) CreateOutboundCall(c *gin.Context) {
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
	conversation.CreateCache(conversationHash, callParams.FirstMessage, callParams.Prompt, callParams.DynamicVariables)
	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"call_sid":          sid,
		"call_status":       status,
		"conversation_hash": conversationHash,
	})
}

func (h *Handler) OutboundCallTwiml(c *gin.Context) {
	conversationHash := c.Query("conversation_hash")
	response := Response{
		Connect: Connect{
			Stream: Stream{
				URL: h.URL.GetRouteUrl("wss", "calls.outbound.call.stream.%hash", conversationHash),
				Parameters: []Parameter{
					{Name: "conversation_hash", Value: conversationHash},
				},
			},
		},
	}

	c.XML(http.StatusOK, response)
}
