package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/internal/services"
	"github.com/starbops/gottodo/ui/templates"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	todoService *services.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// GetAllTodos handles GET /todos
func (h *TodoHandler) GetAllTodos(c echo.Context) error {
	userID := c.Get("user_id").(string)
	todos, err := h.todoService.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, todos)
}

// GetTodo handles GET /todos/:id
func (h *TodoHandler) GetTodo(c echo.Context) error {
	todoID := c.Param("id")
	userID := c.Get("user_id").(string)

	todo, err := h.todoService.GetTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, todo)
}

// CreateTodoRequest represents the request body for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

// UpdateTodoRequest represents the request body for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	// Parse form data
	title := c.FormValue("title")
	description := c.FormValue("description")

	// Validate input
	if title == "" {
		return templates.TodoListWithError("Title is required", nil).Render(c.Request().Context(), c.Response().Writer)
	}

	// Create todo
	todo := &models.Todo{
		Title:       title,
		Description: description,
		UserID:      userID,
		Completed:   false,
	}

	err := h.todoService.CreateTodo(c.Request().Context(), todo)
	if err != nil {
		return templates.TodoListWithError(fmt.Sprintf("Failed to create todo: %v", err), nil).Render(c.Request().Context(), c.Response().Writer)
	}

	// Get updated list of todos
	todos, err := h.todoService.GetUserTodos(c.Request().Context(), userID)
	if err != nil {
		return templates.TodoListWithError(fmt.Sprintf("Failed to retrieve todos: %v", err), nil).Render(c.Request().Context(), c.Response().Writer)
	}

	// Return updated todo list
	return templates.TodoListComponent(todos).Render(c.Request().Context(), c.Response().Writer)
}

// UpdateTodo handles PUT /todos/:id
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	todoID := c.Param("id")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Get user ID and check ownership
	userID := c.Get("user_id").(string)

	// First check if the todo exists and belongs to this user
	_, err := h.todoService.GetTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	todo, err := h.todoService.UpdateTodo(c.Request().Context(), todoID, req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, todo)
}

// UpdateTodoStatus handles PUT /todos/:id/complete and /todos/:id/incomplete
func (h *TodoHandler) UpdateTodoStatus(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	// Get todo ID from URL
	todoID := c.Param("id")

	// Get the current path to determine whether to mark as complete or incomplete
	path := c.Path()
	var completed bool
	if path == "/todos/:id/complete" {
		completed = true
	} else {
		completed = false
	}

	// Update todo
	err := h.todoService.UpdateTodoStatus(c.Request().Context(), todoID, userID, completed)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to update todo: %v", err))
	}

	// Get updated todo
	todo, err := h.todoService.GetTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get updated todo: %v", err))
	}

	// Return updated todo HTML
	return templates.TodoItem(todo).Render(c.Request().Context(), c.Response().Writer)
}

// DeleteTodo handles DELETE /todos/:id
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	// Get user ID from context
	userID := c.Get("user_id").(string)

	// Get todo ID from URL
	todoID := c.Param("id")

	// Delete todo
	err := h.todoService.DeleteTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to delete todo: %v", err))
	}

	// Return empty string to remove the todo from the UI
	return c.NoContent(http.StatusOK)
}
