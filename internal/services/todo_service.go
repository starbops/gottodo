package services

import (
	"context"

	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/internal/repositories"
)

// TodoService handles business logic for todos
type TodoService struct {
	repo repositories.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService(repo repositories.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(ctx context.Context, userID, title, description string) (*models.Todo, error) {
	todo := models.NewTodo(userID, title, description)
	err := s.repo.Create(ctx, todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// GetTodo gets a todo by ID
func (s *TodoService) GetTodo(ctx context.Context, id string) (*models.Todo, error) {
	return s.repo.GetByID(ctx, id)
}

// GetUserTodos gets all todos for a user
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// UpdateTodo updates a todo
func (s *TodoService) UpdateTodo(ctx context.Context, id, title, description string) (*models.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.Update(title, description)
	err = s.repo.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// CompleteTodo marks a todo as completed
func (s *TodoService) CompleteTodo(ctx context.Context, id string) (*models.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.MarkComplete()
	err = s.repo.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// IncompleteTodo marks a todo as not completed
func (s *TodoService) IncompleteTodo(ctx context.Context, id string) (*models.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	todo.MarkIncomplete()
	err = s.repo.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}
