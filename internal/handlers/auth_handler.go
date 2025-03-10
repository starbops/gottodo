package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/pkg/auth"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	service *auth.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	user, err := h.service.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Validate that the user ID is a valid UUID
	if !models.IsValidUUID(user.ID) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Generated user ID is not a valid UUID",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	session, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// Validate that the user ID is a valid UUID
	if !models.IsValidUUID(session.UserID) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Session user ID is not a valid UUID",
		})
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = session.Token
	cookie.Expires = session.ExpiresAt
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":    session.UserID,
		"expires": session.ExpiresAt,
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No authentication token",
		})
	}

	err = h.service.Logout(c.Request().Context(), cookie.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Clear cookie
	cookie = new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return c.NoContent(http.StatusNoContent)
}

// GitHubAuth handles GET /auth/github
func (h *AuthHandler) GitHubAuth(c echo.Context) error {
	// Get the GitHub auth URL and state
	url, state := h.service.GetGitHubAuthURL()

	// Store the state in a cookie for validation
	stateCookie := new(http.Cookie)
	stateCookie.Name = "oauth_state"
	stateCookie.Value = state
	stateCookie.Expires = time.Now().Add(10 * time.Minute)
	stateCookie.Path = "/"
	stateCookie.HttpOnly = true
	stateCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(stateCookie)

	// Redirect to GitHub
	return c.Redirect(http.StatusFound, url)
}

// GitHubCallback handles GET /auth/github/callback
func (h *AuthHandler) GitHubCallback(c echo.Context) error {
	// Get the code and state from the query parameters
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	// Get the state from the cookie
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid OAuth state",
		})
	}

	// Clear the state cookie
	stateCookie = new(http.Cookie)
	stateCookie.Name = "oauth_state"
	stateCookie.Value = ""
	stateCookie.Expires = time.Now().Add(-1 * time.Hour)
	stateCookie.Path = "/"
	stateCookie.HttpOnly = true
	stateCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(stateCookie)

	// Handle the callback
	session, err := h.service.HandleGitHubCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Validate that the user ID is a valid UUID
	if !models.IsValidUUID(session.UserID) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "GitHub auth generated user ID is not a valid UUID",
		})
	}

	// Set the auth cookie
	authCookie := new(http.Cookie)
	authCookie.Name = "auth_token"
	authCookie.Value = session.Token
	authCookie.Expires = session.ExpiresAt
	authCookie.Path = "/"
	authCookie.HttpOnly = true
	authCookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(authCookie)

	// Redirect to the dashboard
	return c.Redirect(http.StatusFound, "/dashboard")
}

// AuthMiddleware is middleware for authenticating requests
func (h *AuthHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authentication required",
			})
		}

		valid, err := h.service.VerifyToken(c.Request().Context(), cookie.Value)
		if err != nil || !valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid authentication token",
			})
		}

		user, err := h.service.GetUser(c.Request().Context(), cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "User not found",
			})
		}

		// Validate that the user ID is a valid UUID
		if !models.IsValidUUID(user.ID) {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "User ID is not a valid UUID",
			})
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)

		return next(c)
	}
}
