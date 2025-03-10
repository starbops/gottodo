package repositories

import (
	"context"
	"sync"

	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/pkg/database"
)

// TodoRepository defines the interface for todo repository operations
type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetByID(ctx context.Context, id string) (*models.Todo, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id string) error
}

// MemoryTodoRepository handles database operations for todos in memory
type MemoryTodoRepository struct {
	db          *database.SupabaseClient
	todos       map[string]*models.Todo   // id -> todo
	todosByUser map[string][]*models.Todo // user_id -> todos
	mu          sync.RWMutex
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(db *database.SupabaseClient) TodoRepository {
	return &MemoryTodoRepository{
		db:          db,
		todos:       make(map[string]*models.Todo),
		todosByUser: make(map[string][]*models.Todo),
	}
}

// Create creates a new todo
func (r *MemoryTodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store the todo
	r.todos[todo.ID] = todo

	// Add to user's todos
	r.todosByUser[todo.UserID] = append(r.todosByUser[todo.UserID], todo)

	return nil
}

// GetByID gets a todo by ID
func (r *MemoryTodoRepository) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

// GetByUserID gets all todos for a user
func (r *MemoryTodoRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos, exists := r.todosByUser[userID]
	if !exists {
		return []*models.Todo{}, nil
	}

	return todos, nil
}

// Update updates a todo
func (r *MemoryTodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return ErrTodoNotFound
	}

	r.todos[todo.ID] = todo

	// Update in user's todos
	for i, t := range r.todosByUser[todo.UserID] {
		if t.ID == todo.ID {
			r.todosByUser[todo.UserID][i] = todo
			break
		}
	}

	return nil
}

// Delete deletes a todo
func (r *MemoryTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	// Remove from user's todos
	userTodos := r.todosByUser[todo.UserID]
	for i, t := range userTodos {
		if t.ID == id {
			r.todosByUser[todo.UserID] = append(userTodos[:i], userTodos[i+1:]...)
			break
		}
	}

	// Remove from todos
	delete(r.todos, id)

	return nil
}
