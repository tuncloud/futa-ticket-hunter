package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const sessionDuration = 7 * 24 * time.Hour

// SessionDuration is the maximum age of a session, exported for use in cookie MaxAge.
const SessionDuration = sessionDuration

type contextKey string

// SessionKey is the context key used to store the current session.
const SessionKey contextKey = "session"

// Session holds data for an authenticated user session.
type Session struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"-"`
}

// SessionFromContext retrieves the Session stored in a context, or nil.
func SessionFromContext(ctx context.Context) *Session {
	s, _ := ctx.Value(SessionKey).(*Session)
	return s
}

// Store is an in-memory, goroutine-safe session store.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewStore creates a new Store and starts a background GC goroutine.
func NewStore() *Store {
	s := &Store{sessions: make(map[string]*Session)}
	go s.gc()
	return s
}

// Create stores a new session for the given user and returns the session token.
func (s *Store) Create(email, name, picture string) string {
	token := NewToken()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = &Session{
		Email:     email,
		Name:      name,
		Picture:   picture,
		CreatedAt: time.Now(),
	}
	return token
}

// Get returns the Session for the given token, or nil if not found or expired.
func (s *Store) Get(token string) *Session {
	s.mu.RLock()
	sess := s.sessions[token]
	s.mu.RUnlock()
	if sess == nil {
		return nil
	}
	if time.Since(sess.CreatedAt) > sessionDuration {
		s.Delete(token)
		return nil
	}
	return sess
}

// Delete removes the session for the given token.
func (s *Store) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

// gc periodically removes expired sessions.
func (s *Store) gc() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for token, sess := range s.sessions {
			if now.Sub(sess.CreatedAt) > sessionDuration {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// SessionStore is the interface for creating, retrieving, and deleting sessions.
type SessionStore interface {
	Create(email, name, picture string) string
	Get(token string) *Session
	Delete(token string)
}

// NewToken generates a cryptographically random session token.
func NewToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("auth: crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GoogleTokenInfo holds the relevant fields returned by Google's tokeninfo endpoint.
type GoogleTokenInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Aud     string `json:"aud"`
}

// VerifyGoogleIDToken calls Google's tokeninfo endpoint to validate an ID token.
// It returns an error if the token is invalid or the audience does not match expectedClientID.
func VerifyGoogleIDToken(idToken, expectedClientID string) (*GoogleTokenInfo, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("tokeninfo request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tokeninfo response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token (status %d)", resp.StatusCode)
	}
	var info GoogleTokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("parse tokeninfo: %w", err)
	}
	if info.Aud != expectedClientID {
		return nil, fmt.Errorf("token audience mismatch")
	}
	if info.Email == "" {
		return nil, fmt.Errorf("missing email in token")
	}
	return &info, nil
}
