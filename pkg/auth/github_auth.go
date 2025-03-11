package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// GitHubOAuthConfig holds the GitHub OAuth configuration
type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
	APIURL       string
}

// NewGitHubOAuthConfig creates a new GitHub OAuth configuration
func NewGitHubOAuthConfig() *GitHubOAuthConfig {
	return &GitHubOAuthConfig{
		// These values will be overridden by the configuration
		// when GetAuthCodeURL is called
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"user:email"},
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		APIURL:       "https://api.github.com",
	}
}

// GetAuthCodeURL returns the URL to redirect the user to for authorization
func (c *GitHubOAuthConfig) GetAuthCodeURL(state string, service *AuthService) string {
	// Use values from configuration if they are set
	clientID, clientSecret, redirectURL := service.config.GetGitHubOAuthConfig()

	if clientID != "" {
		c.ClientID = clientID
	}

	if clientSecret != "" {
		c.ClientSecret = clientSecret
	}

	if redirectURL != "" {
		c.RedirectURL = redirectURL
	}

	// Build URL
	u, err := url.Parse(c.AuthURL)
	if err != nil {
		return ""
	}

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
	return strings.Join(c.Scopes, " ")
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (c *GitHubOAuthConfig) ExchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURL)

	req, err := http.NewRequest("POST", c.TokenURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = data.Encode()

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("error from GitHub: %s", result.Error)
	}

	return result.AccessToken, nil
}

// GetUserInfo retrieves the user information from GitHub
func (c *GitHubOAuthConfig) GetUserInfo(accessToken string) (*GitHubUser, error) {
	// Get user data
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user", c.APIURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	// If the email is not public, get the primary email from the email API
	if user.Email == "" {
		email, err := c.GetPrimaryEmail(accessToken)
		if err != nil {
			return nil, fmt.Errorf("error getting primary email: %w", err)
		}
		user.Email = email
	}

	return &user, nil
}

// GetPrimaryEmail retrieves the primary email address from GitHub
func (c *GitHubOAuthConfig) GetPrimaryEmail(accessToken string) (string, error) {
	// Get emails
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/emails", c.APIURL), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	req.Header.Set("Accept", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	// Find the primary email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary verified email found")
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

// GenerateRandomState generates a random state string for CSRF protection
func GenerateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
