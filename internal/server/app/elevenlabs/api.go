package elevenlabs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type result struct {
	SignedURL string `json:"signed_url"`
}

type InitialConfig struct {
	Type                       string                     `json:"type"`
	DynamicVariables           map[string]any             `json:"dynamic_variables,omitempty"`
	ConversationConfigOverride ConversationConfigOverride `json:"conversation_config_override"`
}

type ConversationConfigOverride struct {
	Agent Agent `json:"agent"`
}

type Agent struct {
	Prompt       Prompt `json:"prompt"`
	FirstMessage string `json:"first_message"`
}

type Prompt struct {
	Prompt string `json:"prompt"`
}

const (
	ElevenLabsWs     = "ElevenLabs"
	TwilioWs         = "Twilio"
	conversationHash = "conversation_hash"
)

const (
	agentId   = "agent_id"
	xiApiKey  = "xi-api-key"
	APIDomain = "api.elevenlabs.io"
)

func (s *Service) GetSignedUrl() (string, error) {

	url := url.URL{
		Scheme: "https",
		Host:   APIDomain,
		Path:   "/v1/convai/conversation/get_signed_url",
	}
	q := url.Query()
	q.Add(agentId, s.App.Config.Elevenlabs.AgentId)
	url.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET",
		url.String(),
		nil,
	)
	req.Header.Add(xiApiKey, s.App.Config.Elevenlabs.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error signed URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get signed URL: %s", resp.Status)
	}

	var res result
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.SignedURL, nil
}

func (s *Service) GenerateElevenLabsConfig(firstMessage, prompt string, dynamicVariables map[string]any) InitialConfig {
	config := InitialConfig{
		Type:             "conversation_initiation_client_data",
		DynamicVariables: dynamicVariables,
		ConversationConfigOverride: ConversationConfigOverride{
			Agent: Agent{
				Prompt: Prompt{
					Prompt: prompt,
				},
				FirstMessage: firstMessage,
			},
		},
	}

	return config
}
