package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
)

// GitHubUser represents a GitHub user
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubOAuthConfig holds GitHub OAuth configuration
type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// NewGitHubOAuthConfig creates a new GitHub OAuth configuration
func NewGitHubOAuthConfig() *GitHubOAuthConfig {
	return &GitHubOAuthConfig{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
	}
}

// GetAuthCodeURL returns the URL for GitHub OAuth authorization
func (c *GitHubOAuthConfig) GetAuthCodeURL(state string) string {
	u, _ := url.Parse(githubAuthorizeURL)
	q := u.Query()
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURL)
	q.Set("scope", c.GetScopeString())
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String()
}

// GetScopeString returns the GitHub OAuth scopes as a space-separated string
func (c *GitHubOAuthConfig) GetScopeString() string {
	var scopeString string
	for i, scope := range c.Scopes {
		if i > 0 {
			scopeString += " "
		}
		scopeString += scope
	}
	return scopeString
}

// ExchangeCodeForToken exchanges an OAuth code for an access token
func (c *GitHubOAuthConfig) ExchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURL)

	req, err := http.NewRequest("POST", githubTokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error exchanging code for token: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("error from GitHub: %s", result.Error)
	}

	return result.AccessToken, nil
}

// GetGitHubUser retrieves the GitHub user information using an access token
func GetGitHubUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", githubUserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding user info: %w", err)
	}

	return &user, nil
}

// CreateSessionFromGitHubUser creates a new session from a GitHub user
func (s *AuthService) CreateSessionFromGitHubUser(gitHubUser *GitHubUser) (*Session, error) {
	// Check if user exists
	var user *User
	for _, u := range s.users {
		if u.Email == gitHubUser.Email {
			user = u
			break
		}
	}

	// Create new user if not exists
	if user == nil {
		user = &User{
			ID:        uuid.New().String(),
			Email:     gitHubUser.Email,
			CreatedAt: time.Now(),
		}
		s.mu.Lock()
		s.users[gitHubUser.Email] = user
		s.mu.Unlock()
	}

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
