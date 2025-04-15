package elevenlabs

import "ai-calls/internal/server/app"

type Service struct {
	App *app.Application
}

func New(app *app.Application) *Service {
	return &Service{App: app}
}
