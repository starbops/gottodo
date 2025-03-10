package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents an authenticated session
type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OAuthState represents a state for OAuth flow
type OAuthState struct {
	State     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// AuthService provides authentication functionality
type AuthService struct {
	users       map[string]*User       // email -> user
	sessions    map[string]*Session    // token -> session
	oauthStates map[string]*OAuthState // state -> OAuthState
	github      *GitHubOAuthConfig
	mu          sync.RWMutex
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{
		users:       make(map[string]*User),
		sessions:    make(map[string]*Session),
		oauthStates: make(map[string]*OAuthState),
		github:      NewGitHubOAuthConfig(),
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.users[email]; exists {
		return nil, errors.New("user already exists")
	}

	// Create a new user with a UUID
	user := &User{
		ID:        uuid.New().String(), // Generate a valid UUID string
		Email:     email,
		CreatedAt: time.Now(),
	}

	// Store the user
	s.users[email] = user

	return user, nil
}

// Login authenticates a user and returns a session
func (s *AuthService) Login(ctx context.Context, email, password string) (*Session, error) {
	s.mu.RLock()
	user, exists := s.users[email]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.New("invalid credentials")
	}

	// In a real implementation, we would validate the password here

	// Create a new session
	session := &Session{
		Token:     uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.sessions[session.Token] = session
	s.mu.Unlock()

	return session, nil
}

// Logout invalidates a session
func (s *AuthService) Logout(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[token]; !exists {
		return errors.New("session not found")
	}

	delete(s.sessions, token)
	return nil
}

// GetUser returns the user associated with a session token
func (s *AuthService) GetUser(ctx context.Context, token string) (*User, error) {
	s.mu.RLock()
	session, exists := s.sessions[token]
	s.mu.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		if exists {
			// Clean up expired session
			s.mu.Lock()
			delete(s.sessions, token)
			s.mu.Unlock()
		}
		return nil, errors.New("invalid or expired session")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find the user by ID
	for _, user := range s.users {
		if user.ID == session.UserID {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// VerifyToken checks if a token is valid
func (s *AuthService) VerifyToken(ctx context.Context, token string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return false, nil
	}

	return true, nil
}

// CreateOAuthState creates a new OAuth state
func (s *AuthService) CreateOAuthState() string {
	state := uuid.New().String()
	now := time.Now()

	oauthState := &OAuthState{
		State:     state,
		CreatedAt: now,
		ExpiresAt: now.Add(10 * time.Minute),
	}

	s.mu.Lock()
	s.oauthStates[state] = oauthState
	s.mu.Unlock()

	return state
}

// VerifyOAuthState verifies an OAuth state and removes it if valid
func (s *AuthService) VerifyOAuthState(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	oauthState, exists := s.oauthStates[state]
	if !exists || time.Now().After(oauthState.ExpiresAt) {
		return false
	}

	// Remove the state after it's used
	delete(s.oauthStates, state)

	return true
}

// GetGitHubAuthURL returns the GitHub OAuth URL
func (s *AuthService) GetGitHubAuthURL() (string, string) {
	state := s.CreateOAuthState()
	return s.github.GetAuthCodeURL(state), state
}

// HandleGitHubCallback handles the GitHub OAuth callback
func (s *AuthService) HandleGitHubCallback(ctx context.Context, code, state string) (*Session, error) {
	// Verify the state
	if !s.VerifyOAuthState(state) {
		return nil, errors.New("invalid OAuth state")
	}

	// Exchange code for token
	accessToken, err := s.github.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get GitHub user
	gitHubUser, err := GetGitHubUser(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	// Create or get user
	return s.CreateSessionFromGitHubUser(gitHubUser)
}
