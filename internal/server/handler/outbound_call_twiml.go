package handler

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) OutboundCallTwiml(c *gin.Context) {
	conversationHash := c.Query("conversation_hash")
	response := Response{
		Connect: Connect{
			Stream: Stream{
				URL: h.URL.GetRouteUrl("wss", "calls.outbound.call.stream"),
				Parameters: []Parameter{
					{Name: "conversation_hash", Value: conversationHash},
				},
			},
		},
	}

	c.XML(http.StatusOK, response)
}
