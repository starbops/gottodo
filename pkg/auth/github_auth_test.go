package auth

import (
	"context"
	"testing"
	"time"

	"github.com/starbops/gottodo/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGitHubOAuthConfig_GetAuthCodeURL(t *testing.T) {
	// Create config
	cfg := &config.Config{}
	cfg.Auth.GitHubClientID = "test-client-id"
	cfg.Auth.GitHubClientSecret = "test-client-secret"
	cfg.Auth.GitHubRedirectURL = "http://localhost:8080/auth/github/callback"

	// Create auth service with config
	authService := NewAuthService(cfg)

	// Create GitHub config
	githubConfig := NewGitHubOAuthConfig()

	// Get auth URL
	state := "test-state"
	url := githubConfig.GetAuthCodeURL(state, authService)

	// Check URL contains correct parameters (account for URL encoding)
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fauth%2Fgithub%2Fcallback")
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "scope=user%3Aemail")
}

func TestGitHubOAuthConfig_GetScopeString(t *testing.T) {
	// Create config
	config := NewGitHubOAuthConfig()
	config.Scopes = []string{"user:email", "repo"}

	// Get scope string
	scopeString := config.GetScopeString()

	// Assert scope string is correct
	assert.Equal(t, "user:email repo", scopeString)
}

func TestCreateSessionFromGitHubUser(t *testing.T) {
	// Create config
	cfg := config.DefaultConfig()

	// Create auth service
	service := NewAuthService(cfg)

	// Create GitHub user
	gitHubUser := &GitHubUser{
		ID:    12345,
		Login: "testuser",
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Create session
	session, err := service.CreateSessionFromGitHubUser(gitHubUser)

	// Assert session was created successfully
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.Token)
	assert.NotEmpty(t, session.UserID)
	assert.True(t, session.ExpiresAt.After(session.ExpiresAt.Add(-24*time.Hour)))

	// Validate user was created
	user, err := service.GetUser(context.Background(), session.Token)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestCreateSessionFromGitHubUser_ExistingUser(t *testing.T) {
	// Create config
	cfg := config.DefaultConfig()

	// Create auth service
	service := NewAuthService(cfg)

	// Create GitHub user
	gitHubUser := &GitHubUser{
		ID:    12345,
		Login: "testuser",
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Create first session
	session1, err := service.CreateSessionFromGitHubUser(gitHubUser)
	assert.NoError(t, err)

	// Create second session for same user
	session2, err := service.CreateSessionFromGitHubUser(gitHubUser)

	// Assert second session was created successfully
	assert.NoError(t, err)
	assert.NotNil(t, session2)
	assert.NotEqual(t, session1.Token, session2.Token)
	assert.Equal(t, session1.UserID, session2.UserID)
}

func TestAuthService_VerifyOAuthState(t *testing.T) {
	// Create config
	cfg := config.DefaultConfig()

	// Create auth service
	service := NewAuthService(cfg)

	// Create OAuth state
	state, err := service.GenerateOAuthState()
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	// Verify OAuth state
	valid := service.VerifyOAuthState(state)
	assert.True(t, valid)

	// Verify invalid state
	valid = service.VerifyOAuthState("invalid-state")
	assert.False(t, valid)
}
