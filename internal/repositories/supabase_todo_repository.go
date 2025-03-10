package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/starbops/gottodo/internal/models"
	"github.com/starbops/gottodo/pkg/database"
)

// SupabaseTodoRepository implements TodoRepository with Supabase PostgreSQL storage
type SupabaseTodoRepository struct {
	db *database.SupabaseClient
}

// NewSupabaseTodoRepository creates a new SupabaseTodoRepository
func NewSupabaseTodoRepository(db *database.SupabaseClient) TodoRepository {
	return &SupabaseTodoRepository{
		db: db,
	}
}

// Create stores a new todo in Supabase
func (r *SupabaseTodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	// Debug: Log the todo object values
	fmt.Printf("Creating todo: ID=%s, UserID=%s, Title='%s', Description='%s', Completed=%v\n",
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Completed)

	// Validate inputs
	if todo.Title == "" {
		return fmt.Errorf("todo title cannot be empty")
	}

	// Parse user ID as UUID
	userID, err := uuid.Parse(todo.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse todo ID as UUID
	todoID, err := uuid.Parse(todo.ID)
	if err != nil {
		return fmt.Errorf("invalid todo ID: %w", err)
	}

	// Debug log the query parameters
	fmt.Printf("SQL Query params: todoID=%v, userID=%v, title='%s', description='%s'\n",
		todoID, userID, todo.Title, todo.Description)

	query := `
		INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Use plain SQL query to debug
	var returnedID uuid.UUID
	err = r.db.DB.QueryRowContext(
		ctx,
		query,
		todoID,           // $1: todo ID (UUID)
		userID,           // $2: user ID (UUID)
		todo.Title,       // $3: title (string)
		todo.Description, // $4: description (string)
		todo.Completed,   // $5: completed (bool)
		todo.CreatedAt,   // $6: created_at (time)
		todo.UpdatedAt,   // $7: updated_at (time)
	).Scan(&returnedID)

	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	// Debug: Log the result
	fmt.Printf("Todo created successfully with ID: %s\n", returnedID)

	// Double-check by fetching the new todo
	newTodo, err := r.GetByID(ctx, returnedID.String())
	if err != nil {
		fmt.Printf("Warning: Todo created but couldn't be retrieved: %v\n", err)
	} else {
		fmt.Printf("Retrieved new todo: Title='%s', Description='%s'\n",
			newTodo.Title, newTodo.Description)
	}

	return nil
}

// GetByID retrieves a todo by ID from Supabase
func (r *SupabaseTodoRepository) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	// Parse todo ID as UUID
	todoID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid todo ID: %w", err)
	}

	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	row := r.db.DB.QueryRowContext(ctx, query, todoID)

	var dbTodoID, dbUserID uuid.UUID
	todo := &models.Todo{}
	err = row.Scan(
		&dbTodoID,
		&dbUserID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTodoNotFound
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	// Convert UUIDs to strings for the model
	todo.ID = dbTodoID.String()
	todo.UserID = dbUserID.String()

	return todo, nil
}

// GetByUserID retrieves all todos for a user from Supabase
func (r *SupabaseTodoRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Todo, error) {
	// Parse user ID as UUID
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.DB.QueryContext(ctx, query, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var dbTodoID, dbUserID uuid.UUID
		todo := &models.Todo{}
		err := rows.Scan(
			&dbTodoID,
			&dbUserID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}

		// Convert UUIDs to strings for the model
		todo.ID = dbTodoID.String()
		todo.UserID = dbUserID.String()

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}

// Update updates a todo in Supabase
func (r *SupabaseTodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	// Parse todo ID as UUID
	todoID, err := uuid.Parse(todo.ID)
	if err != nil {
		return fmt.Errorf("invalid todo ID: %w", err)
	}

	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
	`
	res, err := r.db.DB.ExecContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		time.Now(),
		todoID,
	)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

// Delete deletes a todo from Supabase
func (r *SupabaseTodoRepository) Delete(ctx context.Context, id string) error {
	// Parse todo ID as UUID
	todoID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid todo ID: %w", err)
	}

	query := `DELETE FROM todos WHERE id = $1`
	res, err := r.db.DB.ExecContext(ctx, query, todoID)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
