package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/pkg/auth"
	"github.com/starbops/gottodo/ui/templates"
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
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Error("Register binding error:", err)
		component := templates.RegisterErrorForm("Invalid form data. Please check your inputs.", req.Email)
		return component.Render(context.Background(), c.Response().Writer)
	}

	// Register the user but we don't need the returned user object
	_, err := h.service.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		c.Logger().Error("Registration error:", err)
		component := templates.RegisterErrorForm(err.Error(), req.Email)
		return component.Render(context.Background(), c.Response().Writer)
	}

	// For successful registration, show success message with countdown
	if c.Request().Header.Get("HX-Request") == "true" {
		// Return success component for HTMX requests
		component := templates.RegisterSuccessForm(req.Email)
		return component.Render(context.Background(), c.Response().Writer)
	}

	// For non-HTMX requests, use standard redirect
	return c.Redirect(http.StatusFound, "/login")
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		// Return generic error without revealing details - invalid form data
		return renderLoginError(c, req.Email)
	}

	session, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		// Return generic error without revealing whether the account exists or password is incorrect
		log.Printf("Login failed for email %s: %v", req.Email, err)
		return renderLoginError(c, req.Email)
	}

	// Validate that the user ID is a valid UUID
	if !models.IsValidUUID(session.UserID) {
		// Internal error, but still show generic error to the user
		log.Printf("Invalid UUID for user: %s", session.UserID)
		return renderLoginError(c, req.Email)
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

	// Redirect to dashboard after successful login
	c.Response().Header().Set("HX-Redirect", "/dashboard")
	return c.NoContent(http.StatusOK)
}

// renderLoginError renders a login error form
func renderLoginError(c echo.Context, email string) error {
	return templates.LoginErrorForm(email).Render(c.Request().Context(), c.Response().Writer)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	err = h.service.Logout(c.Request().Context(), cookie.Value)
	if err != nil {
		// Even if there's an error with logout, still clear the cookie
		// and redirect to the login page
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

	// Set both regular redirect and HTMX redirect
	c.Response().Header().Set("HX-Redirect", "/login")

	// Redirect to login page after logout
	return c.Redirect(http.StatusFound, "/login")
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
			// Redirect to login page instead of returning JSON error
			return c.Redirect(http.StatusFound, "/login")
		}

		valid, err := h.service.VerifyToken(c.Request().Context(), cookie.Value)
		if err != nil || !valid {
			// Clear the invalid cookie
			cookie = new(http.Cookie)
			cookie.Name = "auth_token"
			cookie.Value = ""
			cookie.Expires = time.Now().Add(-1 * time.Hour)
			cookie.Path = "/"
			cookie.HttpOnly = true
			cookie.SameSite = http.SameSiteStrictMode
			c.SetCookie(cookie)

			// Redirect to login page instead of returning JSON error
			return c.Redirect(http.StatusFound, "/login")
		}

		user, err := h.service.GetUser(c.Request().Context(), cookie.Value)
		if err != nil {
			// Redirect to login page instead of returning JSON error
			return c.Redirect(http.StatusFound, "/login")
		}

		// Validate that the user ID is a valid UUID
		if !models.IsValidUUID(user.ID) {
			// Redirect to login page instead of returning JSON error
			return c.Redirect(http.StatusFound, "/login")
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)

		return next(c)
	}
}
