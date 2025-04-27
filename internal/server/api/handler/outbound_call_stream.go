package handler

import (
	"ai-calls/internal/logger"
	"ai-calls/internal/server/app/conversations"
	"ai-calls/internal/server/app/elevenlabs"
	"ai-calls/internal/server/app/systemStatus"
	"context"
	"net/http"

	websocket "ai-calls/internal/server/utils/webSocket"

	"github.com/gin-gonic/gin"
)

type ElevenLabsMessage struct {
	Type           string          `json:"type"`
	PingEvent      *PingEvent      `json:"ping_event,omitempty"`
	Audio          *Audio          `json:"audio,omitempty"`
	AudioEvent     *AudioEvent     `json:"audio_event,omitempty"`
	AgentResponse  *AgentResponse  `json:"agent_response_event,omitempty"`
	UserTranscript *UserTranscript `json:"user_transcription_event,omitempty"`
}

type PingEvent struct {
	EventID int `json:"event_id"`
	PingMs  int `json:"ping_ms"`
}

type Audio struct {
	Chunk string `json:"chunk"`
}

type UserAudioChunk struct {
	Chunk string `json:"user_audio_chunk"`
}

type AudioEvent struct {
	AudioBase64 string `json:"audio_base_64"`
}

type AgentResponse struct {
	AgentResponse string `json:"agent_response"`
}

type UserTranscript struct {
	UserTranscript string `json:"user_transcript"`
}

type PongMessage struct {
	Type    string `json:"type"`
	EventID int    `json:"event_id"`
}

type MediaPayload struct {
	Payload string `json:"payload"`
}

type MediaMessage struct {
	Event     string       `json:"event"`
	StreamSid string       `json:"streamSid"`
	Media     MediaPayload `json:"media"`
}

type ClearMessage struct {
	Event     string `json:"event"`
	StreamSid string `json:"streamSid"`
}

type MessageFromTwilio struct {
	Event string `json:"event"`
	Start struct {
		StreamSid        string            `json:"streamSid"`
		CallSid          string            `json:"callSid"`
		CustomParameters map[string]string `json:"customParameters"`
	} `json:"start"`
	Media struct {
		Payload string `json:"payload"`
	} `json:"media"`
}

const (
	ElevenLabsWs     = "ElevenLabs"
	TwilioWs         = "Twilio"
	conversationHash = "conversation_hash"
)

func (h *Handler) OutboundCallStream(c *gin.Context) {
	sysStatus := systemStatus.New(h.App)
	conversation := conversations.New(h.App)
	conversationHash := c.Param("hash")
	firstMessage, prompt, callSidConv, dynamicVariables, err := conversation.GetByHashFromCache(conversationHash)
	if err != nil {
		logger.Error("Conversation failed", logger.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "conversation error",
		})
		return
	}
	conversation.Delete(conversationHash)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	twilioConn, err := websocket.NewServer(h.App.Config, c.Writer, c.Request, nil)
	twilioConn.SetConnectName(TwilioWs)
	if err != nil {
		logger.Error("WebSocket upgrade failed:", logger.Any("err", err))
		return
	}
	defer sysStatus.RemoveWebSocketConnecionCount()
	defer closeConnectionService(ctx, twilioConn)
	sysStatus.AddWebSocketConnecionCount()

	var streamSid, callSid string

	eleven := elevenlabs.New(h.App)
	signedUrl, err := eleven.GetSignedUrl()
	if err != nil {
		cancel()
		logger.Error("[ElevenLabs] get sign url failed:", logger.Any("err", err))
		return
	}

	elevenLabsConn, _, err := websocket.NewClient(ctx, h.App.Config, signedUrl, nil)
	elevenLabsConn.SetConnectName(ElevenLabsWs)
	if err != nil {
		logger.Error("[ElevenLabs] connection error:", logger.Any("err", err))
		return
	}
	logger.Debug("Websocket connected to ElevenLabs")
	defer closeConnectionService(ctx, elevenLabsConn)

	go handleElevenLabsMessages(ctx, elevenLabsConn, twilioConn, &streamSid)

	for {
		var m MessageFromTwilio
		err := twilioConn.ReadJsonMessage(ctx, &m)
		if err != nil {
			logger.Error("[Twilio] Read error:", logger.Any("err", err))
			cancel()
			break
		}

		switch m.Event {
		case "start":
			streamSid = m.Start.StreamSid
			callSid = m.Start.CallSid
			if callSidConv != callSid {
				logger.Error("Wrong Conversation")
				cancel()
				break
			}
			logger.Debug("[Twilio] Stream started", streamSid, callSid)
			elevenlabsInitialConfig := eleven.GenerateElevenLabsConfig(firstMessage, prompt, dynamicVariables)
			err := elevenLabsConn.WriteJsonMessage(ctx, elevenlabsInitialConfig)
			if err != nil {
				logger.Error("Initial Config error:", logger.Any("err", err))
				return
			}

		case "media":
			if elevenLabsConn != nil {
				go func() {
					audio := UserAudioChunk{
						Chunk: m.Media.Payload,
					}
					elevenLabsConn.WriteJsonMessage(ctx, audio)
				}()
			}

		case "stop":
			logger.Debug("[Twilio] Stream stopped", logger.String("id", streamSid))
			cancel()
			return
		}
	}
}

func handleElevenLabsMessages(ctx context.Context, elevenConn *websocket.SafeConn, twilioConn *websocket.SafeConn, streamSid *string) {
	for {
		select {
		case <-ctx.Done():
			logger.Debug("[ElevenLabs] Goroutine stopped by context")
			return
		default:
			var payload ElevenLabsMessage
			err := elevenConn.ReadJsonMessage(ctx, &payload)
			if err != nil {
				logger.Error("[ElevenLabs] Read error:", logger.Any("err", err))
				closeConnectionService(ctx, twilioConn)
				return
			}

			switch payload.Type {
			case "audio":
				if *streamSid != "" {
					go func() {
						media := MediaMessage{
							Event:     "media",
							StreamSid: *streamSid,
							Media: MediaPayload{
								Payload: extractAudioChunk(&payload),
							},
						}
						twilioConn.WriteJsonMessage(ctx, media)
					}()
				}
			case "ping":
				pong := PongMessage{
					Type:    "pong",
					EventID: payload.PingEvent.EventID,
				}
				elevenConn.WriteJsonMessage(ctx, pong)
			case "interruption":
				if *streamSid != "" {
					go func() {
						logger.Debug("[ElevenLabs] interruption")
						media := ClearMessage{
							Event:     "clear",
							StreamSid: *streamSid,
						}
						twilioConn.WriteJsonMessage(ctx, media)
					}()
				}
			case "agent_response":
				logger.Debug("[Agent] response", logger.String("text", payload.AgentResponse.AgentResponse))
			case "user_transcript":
				logger.Debug("[User] response", logger.String("text", payload.UserTranscript.UserTranscript))
			default:
				logger.Error("[ElevenLabs] Ignored type:", logger.String("typeKey", payload.Type))
			}
		}
	}
}

func closeConnectionService(ctx context.Context, conn *websocket.SafeConn) {
	if conn.IsClosed {
		return
	}
	logger.Debug("websocket closed", logger.String("connection", conn.ConnectName()))
	conn.Close(ctx)
}

func extractAudioChunk(msg *ElevenLabsMessage) string {
	if msg.Audio != nil && msg.Audio.Chunk != "" {
		return msg.Audio.Chunk
	}
	if msg.AudioEvent != nil && msg.AudioEvent.AudioBase64 != "" {
		return msg.AudioEvent.AudioBase64
	}
	return ""
}
