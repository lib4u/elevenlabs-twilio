package calls

import (
	"ai-calls/internal/server/app"

	twilio "github.com/twilio/twilio-go"
)

type Service struct {
	App    *app.Application
	Client *twilio.RestClient
}

func New(app *app.Application) *Service {

	c := &Service{App: app}
	c.Client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   c.App.Config.Twilio.APIKey,
		Password:   c.App.Config.Twilio.APISecret,
		AccountSid: c.App.Config.Twilio.AccountSID,
	})
	return c
}
