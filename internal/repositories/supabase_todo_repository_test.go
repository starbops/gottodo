package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/starbops/gottodo/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}

	return db, mock
}

// parseUUID is a helper that safely parses a UUID string and fails the test if invalid
func parseUUID(t *testing.T, uuidStr string) uuid.UUID {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Fatalf("Invalid UUID %s: %v", uuidStr, err)
	}
	return id
}

func TestSupabaseTodoRepository_CreateTodo(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "This is a test todo",
		Completed:   false,
	}

	// Parse userID into UUID for matching in SQL mock
	userUUID := parseUUID(t, userID)

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO todos (id, title, description, user_id, completed) VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(todoID, "Test Todo", "This is a test todo", userUUID, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the function being tested
	err := repo.CreateTodo(ctx, todo)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetTodo(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()

	// Set expected query and response
	rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "completed"}).
		AddRow(todoID, "Test Todo", "This is a test todo", userID, false)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, user_id, completed FROM todos WHERE id = $1`)).
		WithArgs(todoID).
		WillReturnRows(rows)

	// Execute the function being tested
	todo, err := repo.GetTodo(ctx, todoID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, todoID, todo.ID)
	assert.Equal(t, userID, todo.UserID)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.Equal(t, "This is a test todo", todo.Description)
	assert.False(t, todo.Completed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetTodo_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()

	// Set expected query and response for a todo that doesn't exist
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, user_id, completed FROM todos WHERE id = $1`)).
		WithArgs(todoID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function being tested
	_, err := repo.GetTodo(ctx, todoID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetUserTodos(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	userID := uuid.New().String()
	todoID1 := uuid.New().String()
	todoID2 := uuid.New().String()

	// Parse UUIDs for matching in SQL mock
	userUUID := parseUUID(t, userID)

	// Set expected query and response
	rows := sqlmock.NewRows([]string{"id", "title", "description", "user_id", "completed"}).
		AddRow(todoID1, "Todo 1", "Description 1", userID, false).
		AddRow(todoID2, "Todo 2", "Description 2", userID, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, user_id, completed FROM todos WHERE user_id = $1`)).
		WithArgs(userUUID).
		WillReturnRows(rows)

	// Execute the function being tested
	todos, err := repo.GetUserTodos(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, todoID1, todos[0].ID)
	assert.Equal(t, "Todo 1", todos[0].Title)
	assert.Equal(t, todoID2, todos[1].ID)
	assert.Equal(t, "Todo 2", todos[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_UpdateTodo(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Updated Todo",
		Description: "This is an updated test todo",
		Completed:   true,
	}

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE todos SET title = $1, description = $2, completed = $3 WHERE id = $4`)).
		WithArgs("Updated Todo", "This is an updated test todo", true, todoID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the function being tested
	err := repo.UpdateTodo(ctx, todo)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_UpdateTodo_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Updated Todo",
		Description: "This is an updated test todo",
		Completed:   true,
	}

	// Set expected query and response (no rows affected)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE todos SET title = $1, description = $2, completed = $3 WHERE id = $4`)).
		WithArgs("Updated Todo", "This is an updated test todo", true, todoID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute the function being tested
	err := repo.UpdateTodo(ctx, todo)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_DeleteTodo(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM todos WHERE id = $1`)).
		WithArgs(todoID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the function being tested
	err := repo.DeleteTodo(ctx, todoID)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_DeleteTodo_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()

	// Set expected query and response (no rows affected)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM todos WHERE id = $1`)).
		WithArgs(todoID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute the function being tested
	err := repo.DeleteTodo(ctx, todoID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
