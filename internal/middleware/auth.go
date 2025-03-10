package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/pkg/auth"
)

// AuthMiddleware creates middleware that validates the authentication token
func AuthMiddleware(authService *auth.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from cookie
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}

			// Verify token
			valid, err := authService.VerifyToken(c.Request().Context(), cookie.Value)
			if err != nil || !valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authentication token",
				})
			}

			// Get user associated with token
			user, err := authService.GetUser(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User not found",
				})
			}

			// Set user in context
			c.Set("user", user)
			c.Set("user_id", user.ID)

			return next(c)
		}
	}
}
