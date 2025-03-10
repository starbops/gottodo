package repositories

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/starbops/gottodo/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryTodoRepository_GetUserTodos(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	// Create test users and todos
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()

	todo1 := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo 1",
		Description: "Description 1",
		Completed:   false,
		UserID:      userID1,
	}

	todo2 := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo 2",
		Description: "Description 2",
		Completed:   true,
		UserID:      userID1,
	}

	todo3 := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo 3",
		Description: "Description 3",
		Completed:   false,
		UserID:      userID2,
	}

	// Add todos to the repository
	err := repo.CreateTodo(ctx, todo1)
	assert.NoError(t, err)

	err = repo.CreateTodo(ctx, todo2)
	assert.NoError(t, err)

	err = repo.CreateTodo(ctx, todo3)
	assert.NoError(t, err)

	// Test getting todos for user1
	todos, err := repo.GetUserTodos(ctx, userID1)
	assert.NoError(t, err)
	assert.Len(t, todos, 2)

	// Test getting todos for user2
	todos, err = repo.GetUserTodos(ctx, userID2)
	assert.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, todo3.Title, todos[0].Title)

	// Test getting todos for non-existent user
	todos, err = repo.GetUserTodos(ctx, "non-existent-user")
	assert.NoError(t, err)
	assert.Len(t, todos, 0)
}

func TestMemoryTodoRepository_GetTodo(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	// Add todo to the repository
	err := repo.CreateTodo(ctx, todo)
	assert.NoError(t, err)

	// Test getting existing todo
	fetchedTodo, err := repo.GetTodo(ctx, todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, fetchedTodo.ID)
	assert.Equal(t, todo.Title, fetchedTodo.Title)
	assert.Equal(t, todo.Description, fetchedTodo.Description)
	assert.Equal(t, todo.Completed, fetchedTodo.Completed)
	assert.Equal(t, todo.UserID, fetchedTodo.UserID)

	// Test getting non-existent todo
	_, err = repo.GetTodo(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
}

func TestMemoryTodoRepository_CreateTodo(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	userID := uuid.New().String()

	// Test creating a todo with an ID
	todo1 := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo with ID",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	err := repo.CreateTodo(ctx, todo1)
	assert.NoError(t, err)

	fetchedTodo, err := repo.GetTodo(ctx, todo1.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo1.ID, fetchedTodo.ID)

	// Test creating a todo without an ID (should generate one)
	todo2 := &models.Todo{
		Title:       "Test Todo without ID",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	err = repo.CreateTodo(ctx, todo2)
	assert.NoError(t, err)
	assert.NotEmpty(t, todo2.ID)

	fetchedTodo, err = repo.GetTodo(ctx, todo2.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo2.ID, fetchedTodo.ID)
}

func TestMemoryTodoRepository_UpdateTodo(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Original Title",
		Description: "Original Description",
		Completed:   false,
		UserID:      userID,
	}

	// Add todo to the repository
	err := repo.CreateTodo(ctx, todo)
	assert.NoError(t, err)

	// Update the todo
	todo.Title = "Updated Title"
	todo.Description = "Updated Description"
	todo.Completed = true

	err = repo.UpdateTodo(ctx, todo)
	assert.NoError(t, err)

	// Verify the update
	updatedTodo, err := repo.GetTodo(ctx, todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedTodo.Title)
	assert.Equal(t, "Updated Description", updatedTodo.Description)
	assert.True(t, updatedTodo.Completed)

	// Test updating non-existent todo
	nonExistentTodo := &models.Todo{
		ID:          "non-existent-id",
		Title:       "Non-existent Todo",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	err = repo.UpdateTodo(ctx, nonExistentTodo)
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
}

func TestMemoryTodoRepository_DeleteTodo(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          uuid.New().String(),
		Title:       "Test Todo",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	// Add todo to the repository
	err := repo.CreateTodo(ctx, todo)
	assert.NoError(t, err)

	// Delete the todo
	err = repo.DeleteTodo(ctx, todo.ID)
	assert.NoError(t, err)

	// Verify the todo is deleted
	_, err = repo.GetTodo(ctx, todo.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)

	// Test deleting non-existent todo
	err = repo.DeleteTodo(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
}

func TestMemoryTodoRepository_Concurrency(t *testing.T) {
	repo := NewMemoryTodoRepository()
	ctx := context.Background()

	userID := uuid.New().String()
	todoID := uuid.New().String()

	// Create a todo
	todo := &models.Todo{
		ID:          todoID,
		Title:       "Concurrency Test",
		Description: "Description",
		Completed:   false,
		UserID:      userID,
	}

	err := repo.CreateTodo(ctx, todo)
	assert.NoError(t, err)

	// Simulate concurrent operations
	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			_, err := repo.GetTodo(ctx, todoID)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all reads to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify the todo is still intact
	fetchedTodo, err := repo.GetTodo(ctx, todoID)
	assert.NoError(t, err)
	assert.Equal(t, todoID, fetchedTodo.ID)
}
