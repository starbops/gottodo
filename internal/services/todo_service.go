package services

import (
	"context"
	"errors"

	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/internal/repositories"
)

// TodoService handles business logic for todo operations
type TodoService struct {
	todoRepo repositories.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService(todoRepo repositories.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

// GetUserTodos retrieves all todos belonging to a user
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	return s.todoRepo.GetUserTodos(ctx, userID)
}

// GetTodo retrieves a specific todo
func (s *TodoService) GetTodo(ctx context.Context, todoID string, userID string) (*models.Todo, error) {
	if todoID == "" {
		return nil, errors.New("todo ID cannot be empty")
	}

	todo, err := s.todoRepo.GetTodo(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if todo.UserID != userID {
		return nil, errors.New("you don't have permission to access this todo")
	}

	return todo, nil
}

// CreateTodo creates a new todo for a user
func (s *TodoService) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title cannot be empty")
	}

	return s.todoRepo.CreateTodo(ctx, todo)
}

// UpdateTodo updates an existing todo
func (s *TodoService) UpdateTodo(ctx context.Context, todoID string, title string, description string) (*models.Todo, error) {
	// Get the current todo
	todo, err := s.todoRepo.GetTodo(ctx, todoID)
	if err != nil {
		return nil, err
	}

	// Update fields
	todo.Title = title
	todo.Description = description

	// Save changes
	err = s.todoRepo.UpdateTodo(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(ctx context.Context, todoID string, userID string) error {
	// Verify ownership first
	todo, err := s.todoRepo.GetTodo(ctx, todoID)
	if err != nil {
		return err
	}

	if todo.UserID != userID {
		return errors.New("you don't have permission to delete this todo")
	}

	return s.todoRepo.DeleteTodo(ctx, todoID)
}

// UpdateTodoStatus updates the completed status of a todo
func (s *TodoService) UpdateTodoStatus(ctx context.Context, todoID string, userID string, completed bool) error {
	// Verify ownership first
	todo, err := s.todoRepo.GetTodo(ctx, todoID)
	if err != nil {
		return err
	}

	if todo.UserID != userID {
		return errors.New("you don't have permission to update this todo")
	}

	// Update status
	todo.Completed = completed

	// Save changes
	return s.todoRepo.UpdateTodo(ctx, todo)
}
