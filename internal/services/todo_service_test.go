package services

import (
	"context"
	"testing"

	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/internal/repositories"
)

// MockTodoRepository is a mock implementation of the TodoRepository interface
type MockTodoRepository struct {
	todos map[string]*models.Todo
}

// NewMockTodoRepository creates a new MockTodoRepository
func NewMockTodoRepository() repositories.TodoRepository {
	return &MockTodoRepository{
		todos: make(map[string]*models.Todo),
	}
}

// Create implements the Create method of the TodoRepository interface
func (r *MockTodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	r.todos[todo.ID] = todo
	return nil
}

// GetByID implements the GetByID method of the TodoRepository interface
func (r *MockTodoRepository) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return nil, repositories.ErrTodoNotFound
	}
	return todo, nil
}

// GetByUserID implements the GetByUserID method of the TodoRepository interface
func (r *MockTodoRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Todo, error) {
	var todos []*models.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

// Update implements the Update method of the TodoRepository interface
func (r *MockTodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	if _, ok := r.todos[todo.ID]; !ok {
		return repositories.ErrTodoNotFound
	}
	r.todos[todo.ID] = todo
	return nil
}

// Delete implements the Delete method of the TodoRepository interface
func (r *MockTodoRepository) Delete(ctx context.Context, id string) error {
	if _, ok := r.todos[id]; !ok {
		return repositories.ErrTodoNotFound
	}
	delete(r.todos, id)
	return nil
}

func TestCreateTodo(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create a todo
	todo, err := service.CreateTodo(context.Background(), "user123", "Test Todo", "This is a test todo")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Check that the todo was created with the correct values
	if todo.Title != "Test Todo" {
		t.Errorf("Expected title to be 'Test Todo', got '%s'", todo.Title)
	}
	if todo.Description != "This is a test todo" {
		t.Errorf("Expected description to be 'This is a test todo', got '%s'", todo.Description)
	}
	if todo.UserID != "user123" {
		t.Errorf("Expected user ID to be 'user123', got '%s'", todo.UserID)
	}
	if todo.Completed {
		t.Errorf("Expected todo to be incomplete, but it was marked as completed")
	}
}

func TestGetUserTodos(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create some todos for different users
	service.CreateTodo(context.Background(), "user1", "User 1 Todo 1", "Description 1")
	service.CreateTodo(context.Background(), "user1", "User 1 Todo 2", "Description 2")
	service.CreateTodo(context.Background(), "user2", "User 2 Todo", "Description 3")

	// Get todos for user1
	todos, err := service.GetUserTodos(context.Background(), "user1")
	if err != nil {
		t.Fatalf("Failed to get user todos: %v", err)
	}

	// Check that we got the correct number of todos
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos for user1, got %d", len(todos))
	}

	// Get todos for user2
	todos, err = service.GetUserTodos(context.Background(), "user2")
	if err != nil {
		t.Fatalf("Failed to get user todos: %v", err)
	}

	// Check that we got the correct number of todos
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo for user2, got %d", len(todos))
	}
}

func TestCompleteTodo(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create a todo
	todo, err := service.CreateTodo(context.Background(), "user1", "Test Todo", "This is a test todo")
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Complete the todo
	completedTodo, err := service.CompleteTodo(context.Background(), todo.ID)
	if err != nil {
		t.Fatalf("Failed to complete todo: %v", err)
	}

	// Check that the todo was marked as completed
	if !completedTodo.Completed {
		t.Errorf("Expected todo to be completed, but it was not")
	}
}
