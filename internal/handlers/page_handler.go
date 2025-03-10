package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/pkg/auth"
	"github.com/starbops/gottodo/ui/templates"
)

// PageHandler handles HTTP requests for HTML pages
type PageHandler struct {
	todoService *services.TodoService
	authService *auth.AuthService
}

// NewPageHandler creates a new PageHandler
func NewPageHandler(todoService *services.TodoService, authService *auth.AuthService) *PageHandler {
	return &PageHandler{
		todoService: todoService,
		authService: authService,
	}
}

// Home handles GET /
func (h *PageHandler) Home(c echo.Context) error {
	return templates.Home().Render(c.Request().Context(), c.Response().Writer)
}

// Login handles GET /login
func (h *PageHandler) Login(c echo.Context) error {
	return templates.Login().Render(c.Request().Context(), c.Response().Writer)
}

// Register handles GET /register
func (h *PageHandler) Register(c echo.Context) error {
	return templates.Register().Render(c.Request().Context(), c.Response().Writer)
}

// Dashboard handles GET /dashboard
func (h *PageHandler) Dashboard(c echo.Context) error {
	// Get user from context
	userID := c.Get("user_id").(string)

	// Get todos for the user
	todos, err := h.todoService.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Render the dashboard template with the todos
	return templates.Dashboard(todos).Render(c.Request().Context(), c.Response().Writer)
}
