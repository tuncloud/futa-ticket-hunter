package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// ClerkVerifier validates Clerk session tokens using the configured JWKS.
type ClerkVerifier struct {
	issuer string
	jwks   *keyfunc.JWKS
}

// NewClerkVerifier creates a new ClerkVerifier using the issuer or JWKS URL.
func NewClerkVerifier(issuer, jwksURL string) (*ClerkVerifier, error) {
	issuer = strings.TrimRight(strings.TrimSpace(issuer), "/")
	jwksURL = strings.TrimSpace(jwksURL)
	if jwksURL == "" && issuer != "" {
		jwksURL = issuer + "/.well-known/jwks.json"
	}
	if jwksURL == "" {
		return nil, errors.New("missing clerk issuer or jwks url")
	}
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return nil, fmt.Errorf("load jwks: %w", err)
	}
	return &ClerkVerifier{issuer: issuer, jwks: jwks}, nil
}

// VerifySession verifies a Clerk JWT and returns the session payload.
func (v *ClerkVerifier) VerifySession(token string) (*Session, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}
	claims := jwt.MapClaims{}
	options := []jwt.ParserOption{jwt.WithValidMethods([]string{"RS256"})}
	parsed, err := jwt.ParseWithClaims(token, claims, v.jwks.Keyfunc, options...)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	if v.issuer != "" {
		issuer, _ := claims["iss"].(string)
		if strings.TrimRight(issuer, "/") != v.issuer {
			return nil, errors.New("issuer mismatch")
		}
	}
	email := claimString(claims, "email")
	if email == "" {
		email = claimString(claims, "email_address")
	}
	if email == "" {
		email = emailFromAddresses(claims["email_addresses"])
	}
	if email == "" {
		return nil, errors.New("missing email")
	}
	name := claimString(claims, "name")
	if name == "" {
		name = claimString(claims, "full_name")
	}
	if name == "" {
		first := claimString(claims, "first_name")
		last := claimString(claims, "last_name")
		name = strings.TrimSpace(strings.TrimSpace(first + " " + last))
	}
	picture := claimString(claims, "image_url")
	if picture == "" {
		picture = claimString(claims, "picture")
	}
	return &Session{Email: email, Name: name, Picture: picture, CreatedAt: time.Now()}, nil
}

func claimString(claims jwt.MapClaims, key string) string {
	if value, ok := claims[key].(string); ok {
		return value
	}
	return ""
}

func emailFromAddresses(value any) string {
	list, ok := value.([]any)
	if !ok {
		return ""
	}
	for _, item := range list {
		if entry, ok := item.(map[string]any); ok {
			if email, ok := entry["email_address"].(string); ok && email != "" {
				return email
			}
		}
	}
	return ""
}
