package config

import (
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Worker   WorkerConfig   `yaml:"worker"`
	Webhook  WebhookConfig  `yaml:"webhook"`
	Email    EmailConfig    `yaml:"email"`
	Futa     FutaConfig     `yaml:"futa"`
	Clerk    ClerkConfig    `yaml:"clerk"`
	Posthog  PosthogConfig  `yaml:"posthog"`
}

type PosthogConfig struct {
	APIKey string `yaml:"api_key" envconfig:"POSTHOG_API_KEY"`
	Host   string `yaml:"host" envconfig:"POSTHOG_HOST"`
}

type ClerkConfig struct {
	PublishableKey string `yaml:"publishable_key" envconfig:"CLERK_PUBLISHABLE_KEY"`
	SecretKey      string `yaml:"secret_key" envconfig:"CLERK_SECRET_KEY"`
	Issuer         string `yaml:"issuer" envconfig:"CLERK_ISSUER"`
	JWKSURL        string `yaml:"jwks_url" envconfig:"CLERK_JWKS_URL"`
}

type ServerConfig struct {
	Port int `yaml:"port" envconfig:"SERVER_PORT"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" envconfig:"DATABASE_HOST"`
	Port     int    `yaml:"port" envconfig:"DATABASE_PORT"`
	User     string `yaml:"user" envconfig:"DATABASE_USER"`
	Password string `yaml:"password" envconfig:"DATABASE_PASSWORD"`
	DBName   string `yaml:"dbname" envconfig:"DATABASE_DBNAME"`
	SSLMode  string `yaml:"sslmode" envconfig:"DATABASE_SSLMODE"`
}

type WorkerConfig struct {
	PollInterval time.Duration `yaml:"poll_interval" envconfig:"WORKER_POLL_INTERVAL"`
	MaxRetries   int           `yaml:"max_retries" envconfig:"WORKER_MAX_RETRIES"`
	Concurrency  int           `yaml:"concurrency" envconfig:"WORKER_CONCURRENCY"`
	RetryDelay   time.Duration `yaml:"retry_delay" envconfig:"WORKER_RETRY_DELAY"`
}

type WebhookConfig struct {
	URL    string `yaml:"url" envconfig:"WEBHOOK_URL"`
	Secret string `yaml:"secret" envconfig:"WEBHOOK_SECRET"`
}

type EmailConfig struct {
	ResendAPIKey string `yaml:"resend_api_key" envconfig:"EMAIL_RESEND_API_KEY"`
	FromAddress  string `yaml:"from_address" envconfig:"EMAIL_FROM_ADDRESS"`
	FromName     string `yaml:"from_name" envconfig:"EMAIL_FROM_NAME"`
}

type FutaConfig struct {
	BaseURL    string `yaml:"base_url" envconfig:"FUTA_BASE_URL"`
	WebURL     string `yaml:"web_url" envconfig:"FUTA_WEB_URL"`
	UserAgent  string `yaml:"user_agent" envconfig:"FUTA_USER_AGENT"`
	AppVersion string `yaml:"app_version" envconfig:"FUTA_APP_VERSION"`
	Channel    string `yaml:"channel" envconfig:"FUTA_CHANNEL"`
}

// Load loads config from file and overrides with environment variables using envconfig
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

	// Override bằng envconfig
	// Cần import "github.com/kelseyhightower/envconfig"
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
