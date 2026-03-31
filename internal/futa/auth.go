package futa

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const firebaseTokenURL = "https://securetoken.googleapis.com/v1/token?key=AIzaSyCl2ZuHKk41TcLv5n9_5coBKDttr6PCo-Q"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

type IDTokenClaims struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	CustomUID   int    `json:"custom_uid"`
	Email       string `json:"email"`
}

// ExchangeRefreshToken exchanges a Firebase refresh token for access + id tokens.
func ExchangeRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	body := fmt.Sprintf(`{"grantType":"refresh_token","refreshToken":"%s"}`, refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", firebaseTokenURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("exchange refresh token: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(data))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

// ParseIDToken extracts claims from a Firebase ID token (JWT) without verification.
// This is safe because the token comes directly from Google's token endpoint.
func ParseIDToken(idToken string) (*IDTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	payload := parts[1]
	// Add padding if needed
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("decode JWT payload: %w", err)
	}

	var claims IDTokenClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("parse JWT claims: %w", err)
	}
	return &claims, nil
}
