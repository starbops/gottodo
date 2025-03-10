package repositories

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/starbops/gottodo/internal/models"
)

// MemoryTodoRepository is an in-memory implementation of TodoRepository
type MemoryTodoRepository struct {
	todos map[string]*models.Todo
	mutex sync.RWMutex
}

// NewMemoryTodoRepository creates a new MemoryTodoRepository
func NewMemoryTodoRepository() TodoRepository {
	return &MemoryTodoRepository{
		todos: make(map[string]*models.Todo),
	}
}

// generateID creates a unique ID for a new todo
func generateID() string {
	return uuid.New().String()
}

// GetUserTodos retrieves all todos for a specific user
func (r *MemoryTodoRepository) GetUserTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var userTodos []*models.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}

	return userTodos, nil
}

// GetTodo retrieves a specific todo by ID
func (r *MemoryTodoRepository) GetTodo(ctx context.Context, todoID string) (*models.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todo, exists := r.todos[todoID]
	if !exists {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

// CreateTodo creates a new todo
func (r *MemoryTodoRepository) CreateTodo(ctx context.Context, todo *models.Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Ensure the todo has an ID
	if todo.ID == "" {
		todo.ID = generateID()
	}

	r.todos[todo.ID] = todo
	return nil
}

// UpdateTodo updates an existing todo
func (r *MemoryTodoRepository) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.todos[todo.ID]
	if !exists {
		return ErrTodoNotFound
	}

	r.todos[todo.ID] = todo
	return nil
}

// DeleteTodo deletes a todo by ID
func (r *MemoryTodoRepository) DeleteTodo(ctx context.Context, todoID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.todos[todoID]
	if !exists {
		return ErrTodoNotFound
	}

	delete(r.todos, todoID)
	return nil
}
