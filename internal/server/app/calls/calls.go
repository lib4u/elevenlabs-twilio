package calls

import (
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (c *Service) CallParams() *twilioApi.CreateCallParams {
	return &twilioApi.CreateCallParams{}
}

func (c *Service) InitCall(params *twilioApi.CreateCallParams) (*twilioApi.ApiV2010Call, error) {
	return c.Client.Api.CreateCall(params)
}

func (c *Service) CreateCall(to string, from string, url string) (string, string, error) {
	params := c.CallParams()
	params.SetMachineDetection("Enable")
	params.SetTo(to)
	params.SetFrom(from)
	params.SetUrl(url)
	resp, err := c.InitCall(params)
	if err != nil {
		return "", "", err
	}
	return *resp.Status, *resp.Sid, nil
}
