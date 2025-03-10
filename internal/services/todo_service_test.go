package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

// GetUserTodos implements the GetUserTodos method of the TodoRepository interface
func (r *MockTodoRepository) GetUserTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	var todos []*models.Todo
	for _, todo := range r.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

// GetTodo implements the GetTodo method of the TodoRepository interface
func (r *MockTodoRepository) GetTodo(ctx context.Context, todoID string) (*models.Todo, error) {
	todo, ok := r.todos[todoID]
	if !ok {
		return nil, repositories.ErrTodoNotFound
	}
	return todo, nil
}

// CreateTodo implements the CreateTodo method of the TodoRepository interface
func (r *MockTodoRepository) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.ID == "" {
		todo.ID = uuid.New().String()
	}
	r.todos[todo.ID] = todo
	return nil
}

// UpdateTodo implements the UpdateTodo method of the TodoRepository interface
func (r *MockTodoRepository) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	if _, ok := r.todos[todo.ID]; !ok {
		return repositories.ErrTodoNotFound
	}
	r.todos[todo.ID] = todo
	return nil
}

// DeleteTodo implements the DeleteTodo method of the TodoRepository interface
func (r *MockTodoRepository) DeleteTodo(ctx context.Context, todoID string) error {
	if _, ok := r.todos[todoID]; !ok {
		return repositories.ErrTodoNotFound
	}
	delete(r.todos, todoID)
	return nil
}

func TestCreateTodo(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create a todo
	todo := &models.Todo{
		UserID:      "user123",
		Title:       "Test Todo",
		Description: "This is a test todo",
	}

	err := service.CreateTodo(context.Background(), todo)
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
	if todo.ID == "" {
		t.Errorf("Expected todo ID to be generated, but it was empty")
	}
}

func TestGetUserTodos(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create some todos for different users
	err := service.CreateTodo(context.Background(), &models.Todo{
		UserID:      "user1",
		Title:       "User 1 Todo 1",
		Description: "Description 1",
	})
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	err = service.CreateTodo(context.Background(), &models.Todo{
		UserID:      "user1",
		Title:       "User 1 Todo 2",
		Description: "Description 2",
	})
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	err = service.CreateTodo(context.Background(), &models.Todo{
		UserID:      "user2",
		Title:       "User 2 Todo",
		Description: "Description 3",
	})
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

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

func TestUpdateTodoStatus(t *testing.T) {
	// Create a mock repository
	repo := NewMockTodoRepository()

	// Create a service with the mock repository
	service := NewTodoService(repo)

	// Create a todo
	todo := &models.Todo{
		UserID:      "user1",
		Title:       "Test Todo",
		Description: "This is a test todo",
	}

	err := service.CreateTodo(context.Background(), todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Update the todo status to completed
	err = service.UpdateTodoStatus(context.Background(), todo.ID, "user1", true)
	if err != nil {
		t.Fatalf("Failed to update todo status: %v", err)
	}

	// Get the updated todo
	updatedTodo, err := service.GetTodo(context.Background(), todo.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to get todo: %v", err)
	}

	// Check that the todo was marked as completed
	if !updatedTodo.Completed {
		t.Errorf("Expected todo to be completed, but it was not")
	}
}
