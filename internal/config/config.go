package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Worker   WorkerConfig   `yaml:"worker"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Email    EmailConfig    `yaml:"email"`
	Futa     FutaConfig     `yaml:"futa"`
	Google   GoogleConfig   `yaml:"google"`
	Posthog  PosthogConfig  `yaml:"posthog"`
}

type PosthogConfig struct {
	APIKey string `yaml:"api_key"`
	Host   string `yaml:"host"`
}

type GoogleConfig struct {
	ClientID string `yaml:"client_id"`
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
	Concurrency  int           `yaml:"concurrency"`
	RetryDelay   time.Duration `yaml:"retry_delay"`
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
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Đọc file config
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal vào struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
