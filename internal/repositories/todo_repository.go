package repositories

import (
	"context"

	"github.com/starbops/gottodo/internal/models"
)

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	// GetUserTodos retrieves all todos for a specific user
	GetUserTodos(ctx context.Context, userID string) ([]*models.Todo, error)

	// GetTodo retrieves a specific todo by ID
	GetTodo(ctx context.Context, todoID string) (*models.Todo, error)

	// CreateTodo creates a new todo
	CreateTodo(ctx context.Context, todo *models.Todo) error

	// UpdateTodo updates an existing todo
	UpdateTodo(ctx context.Context, todo *models.Todo) error

	// DeleteTodo deletes a todo by ID
	DeleteTodo(ctx context.Context, todoID string) error
}
