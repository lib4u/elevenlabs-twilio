package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" envDefault:"development"`
	HTTPServer `yaml:"http_server"`
	WebSocket  `yaml:"web_socket"`
	Twilio     `yaml:"twilio"`
	Elevenlabs `yaml:"elevenlabs"`
}

type HTTPServer struct {
	Port           int      `yaml:"port"`
	Host           string   `yaml:"host"`
	TrustedProxies []string `yaml:"trusted_proxies"`
}

type WebSocket struct {
	ReadLimit int64 `yaml:"read_limit"`
}

type Twilio struct {
	AccountSID string `yaml:"account_sid" env:"TWILIO_ACCOUNT_SID"`
	APIKey     string `yaml:"api_key" env:"TWILIO_API_KEY"`
	APISecret  string `yaml:"api_secret" env:"TWILIO_API_SECRET"`
	AppSID     string `yaml:"app_sid" env:"TWILIO_APP_SID"`
	AuthToken  string `yaml:"auth_token" env:"TWILIO_AUTH_TOKEN"`
}

type Elevenlabs struct {
	APIKey  string `yaml:"api_key"`
	AgentId string `yaml:"agent_id"`
}

func Load(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	return cfg
}
