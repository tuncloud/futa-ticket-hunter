package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Worker   WorkerConfig   `yaml:"worker"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Email    EmailConfig    `yaml:"email"`
	Futa     FutaConfig     `yaml:"futa"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type WorkerConfig struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	MaxRetries   int           `yaml:"max_retries"`
}

type WebhookConfig struct {
	URL    string `yaml:"url"`
	Secret string `yaml:"secret"`
}

type EmailConfig struct {
	ResendAPIKey string `yaml:"resend_api_key"`
	FromAddress  string `yaml:"from_address"`
	FromName     string `yaml:"from_name"`
}

type FutaConfig struct {
	BaseURL    string `yaml:"base_url"`
	WebURL     string `yaml:"web_url"`
	UserAgent  string `yaml:"user_agent"`
	AppVersion string `yaml:"app_version"`
	Channel    string `yaml:"channel"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Worker.PollInterval == 0 {
		cfg.Worker.PollInterval = 30 * time.Second
	}
	if cfg.Worker.MaxRetries == 0 {
		cfg.Worker.MaxRetries = 3
	}
	return &cfg, nil
}
