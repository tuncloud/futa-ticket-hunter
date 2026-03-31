package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/database"
)

type Sender struct {
	cfg    config.WebhookConfig
	client *http.Client
}

func NewSender(cfg config.WebhookConfig) *Sender {
	return &Sender{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type Payload struct {
	Event     string                   `json:"event"`
	Timestamp string                   `json:"timestamp"`
	Data      database.BookingSchedule `json:"data"`
}

func (s *Sender) Send(ctx context.Context, schedule database.BookingSchedule) error {
	if s.cfg.URL == "" {
		return nil
	}

	payload := Payload{
		Event:     "booking.success",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      schedule,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if s.cfg.Secret != "" {
		mac := hmac.New(sha256.New, []byte(s.cfg.Secret))
		mac.Write(body)
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Webhook-Signature", sig)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}
