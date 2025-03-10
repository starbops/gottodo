package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item
type Todo struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTodo creates a new Todo item
func NewTodo(userID, title, description string) *Todo {
	now := time.Now()

	// Generate a new UUID for the todo
	todoID := uuid.New().String()

	return &Todo{
		ID:          todoID,
		UserID:      userID, // UserID should be a valid UUID string
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkComplete marks a todo as completed
func (t *Todo) MarkComplete() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkIncomplete marks a todo as not completed
func (t *Todo) MarkIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}

// Update updates the todo's details
func (t *Todo) Update(title, description string) {
	t.Title = title
	t.Description = description
	t.UpdatedAt = time.Now()
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
