package elevenlabs

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type result struct {
	SignedURL string `json:"signed_url"`
}

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
		return "", err
	}
	defer resp.Body.Close()
	var res result
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.SignedURL, nil
}
