package repositories

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/pkg/database"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*database.SupabaseClient, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}

	return &database.SupabaseClient{DB: db}, mock
}

// parseUUID is a helper that safely parses a UUID string and fails the test if invalid
func parseUUID(t *testing.T, uuidStr string) uuid.UUID {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Fatalf("Invalid UUID %s: %v", uuidStr, err)
	}
	return id
}

func TestSupabaseTodoRepository_Create(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	now := time.Now()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Test Todo",
		Description: "This is a test todo",
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Parse UUIDs for matching in SQL mock
	todoUUID := parseUUID(t, todoID)
	userUUID := parseUUID(t, userID)

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)).WithArgs(
		todoUUID,
		userUUID,
		"Test Todo",
		"This is a test todo",
		false,
		now,
		now,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the function being tested
	err := repo.Create(ctx, todo)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetByID(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	now := time.Now()

	// Parse UUIDs for matching in SQL mock
	todoUUID := parseUUID(t, todoID)
	userUUID := parseUUID(t, userID)

	// Set expected query and response
	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "completed", "created_at", "updated_at"}).
		AddRow(todoUUID, userUUID, "Test Todo", "This is a test todo", false, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`)).WithArgs(todoUUID).WillReturnRows(rows)

	// Execute the function being tested
	todo, err := repo.GetByID(ctx, todoID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, todoID, todo.ID)
	assert.Equal(t, userID, todo.UserID)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.Equal(t, "This is a test todo", todo.Description)
	assert.False(t, todo.Completed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetByID_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()
	todoUUID := parseUUID(t, todoID)

	// Set expected query and response for a todo that doesn't exist
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`)).WithArgs(todoUUID).WillReturnError(sql.ErrNoRows)

	// Execute the function being tested
	_, err := repo.GetByID(ctx, todoID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_GetByUserID(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	userID := uuid.New().String()
	todoID1 := uuid.New().String()
	todoID2 := uuid.New().String()
	now := time.Now()

	// Parse UUIDs for matching in SQL mock
	userUUID := parseUUID(t, userID)
	todoUUID1 := parseUUID(t, todoID1)
	todoUUID2 := parseUUID(t, todoID2)

	// Set expected query and response
	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "completed", "created_at", "updated_at"}).
		AddRow(todoUUID1, userUUID, "Todo 1", "Description 1", false, now, now).
		AddRow(todoUUID2, userUUID, "Todo 2", "Description 2", true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`)).WithArgs(userUUID).WillReturnRows(rows)

	// Execute the function being tested
	todos, err := repo.GetByUserID(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, todoID1, todos[0].ID)
	assert.Equal(t, "Todo 1", todos[0].Title)
	assert.Equal(t, todoID2, todos[1].ID)
	assert.Equal(t, "Todo 2", todos[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_Update(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	now := time.Now()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Updated Todo",
		Description: "This is an updated test todo",
		Completed:   true,
		UpdatedAt:   now,
	}

	// Parse UUID for matching in SQL mock
	todoUUID := parseUUID(t, todoID)

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE todos
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
	`)).WithArgs(
		"Updated Todo",
		"This is an updated test todo",
		true,
		sqlmock.AnyArg(),
		todoUUID,
	).WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the function being tested
	err := repo.Update(ctx, todo)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_Update_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create valid UUIDs for testing
	todoID := uuid.New().String()
	userID := uuid.New().String()
	now := time.Now()
	todo := &models.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       "Updated Todo",
		Description: "This is an updated test todo",
		Completed:   true,
		UpdatedAt:   now,
	}

	// Parse UUID for matching in SQL mock
	todoUUID := parseUUID(t, todoID)

	// Set expected query and response for a todo that doesn't exist
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE todos
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
	`)).WithArgs(
		"Updated Todo",
		"This is an updated test todo",
		true,
		sqlmock.AnyArg(),
		todoUUID,
	).WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute the function being tested
	err := repo.Update(ctx, todo)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_Delete(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()
	todoUUID := parseUUID(t, todoID)

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM todos WHERE id = $1`)).
		WithArgs(todoUUID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the function being tested
	err := repo.Delete(ctx, todoID)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupabaseTodoRepository_Delete_NotFound(t *testing.T) {
	// Setup
	mockDB, mock := setupMockDB(t)
	repo := NewSupabaseTodoRepository(mockDB)
	ctx := context.Background()

	// Create a valid UUID for testing
	todoID := uuid.New().String()
	todoUUID := parseUUID(t, todoID)

	// Set expected query and response
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM todos WHERE id = $1`)).
		WithArgs(todoUUID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute the function being tested
	err := repo.Delete(ctx, todoID)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrTodoNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
