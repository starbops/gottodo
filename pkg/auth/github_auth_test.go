package auth

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubOAuthConfig_GetAuthCodeURL(t *testing.T) {
	// Save original environment variables
	originalClientID := os.Getenv("GITHUB_CLIENT_ID")
	originalClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	originalRedirectURL := os.Getenv("GITHUB_REDIRECT_URL")

	// Set test environment variables
	os.Setenv("GITHUB_CLIENT_ID", "test-client-id")
	os.Setenv("GITHUB_CLIENT_SECRET", "test-client-secret")
	os.Setenv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback")

	// Restore original environment variables after test
	defer func() {
		os.Setenv("GITHUB_CLIENT_ID", originalClientID)
		os.Setenv("GITHUB_CLIENT_SECRET", originalClientSecret)
		os.Setenv("GITHUB_REDIRECT_URL", originalRedirectURL)
	}()

	// Create config
	config := NewGitHubOAuthConfig()

	// Generate auth URL
	state := "test-state"
	url := config.GetAuthCodeURL(state)

	// Assert URL contains expected parameters
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fauth%2Fgithub%2Fcallback")
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "scope=user%3Aemail")
}

func TestGitHubOAuthConfig_GetScopeString(t *testing.T) {
	// Create config with specific scopes
	config := &GitHubOAuthConfig{
		Scopes: []string{"user", "repo", "gist"},
	}

	// Get scope string
	scopeString := config.GetScopeString()

	// Assert scope string is as expected
	assert.Equal(t, "user repo gist", scopeString)
}

func TestCreateSessionFromGitHubUser(t *testing.T) {
	// Create auth service
	service := NewAuthService()

	// Create a GitHub user
	gitHubUser := &GitHubUser{
		ID:        12345,
		Login:     "testuser",
		Name:      "Test User",
		Email:     "test@example.com",
		AvatarURL: "https://github.com/avatar.png",
	}

	// Create session
	session, err := service.CreateSessionFromGitHubUser(gitHubUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.Token)
	assert.NotEmpty(t, session.UserID)
	assert.False(t, session.ExpiresAt.IsZero())

	// Verify user was created
	user, err := service.GetUser(context.Background(), session.Token)
	assert.NoError(t, err)
	assert.Equal(t, gitHubUser.Email, user.Email)
}

func TestCreateSessionFromGitHubUser_ExistingUser(t *testing.T) {
	// Create auth service
	service := NewAuthService()

	// Create a user first
	email := "existing@example.com"
	existingUser, err := service.Register(context.Background(), email, "password")
	assert.NoError(t, err)

	// Create a GitHub user with the same email
	gitHubUser := &GitHubUser{
		ID:        12345,
		Login:     "existinguser",
		Name:      "Existing User",
		Email:     email,
		AvatarURL: "https://github.com/avatar.png",
	}

	// Create session
	session, err := service.CreateSessionFromGitHubUser(gitHubUser)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, existingUser.ID, session.UserID) // Should be the same user ID
}

func TestAuthService_VerifyOAuthState(t *testing.T) {
	// Create auth service
	service := NewAuthService()

	// Create a state
	state := service.CreateOAuthState()
	assert.NotEmpty(t, state)

	// Verify the state
	valid := service.VerifyOAuthState(state)
	assert.True(t, valid)

	// Verify the state can't be used again
	valid = service.VerifyOAuthState(state)
	assert.False(t, valid)

	// Verify an invalid state returns false
	valid = service.VerifyOAuthState("invalid-state")
	assert.False(t, valid)
}
