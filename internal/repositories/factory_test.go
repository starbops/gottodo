package repositories

import (
	"testing"

	"github.com/starbops/gottodo/pkg/config"
)

func TestNewTodoRepository_Memory(t *testing.T) {
	// Create default config (which uses memory repository)
	cfg := config.DefaultConfig()

	// Create the repository
	repo, err := NewTodoRepository(cfg)
	if err != nil {
		t.Fatalf("Failed to create memory repository: %v", err)
	}

	// Verify the type of repository
	_, ok := repo.(*MemoryTodoRepository)
	if !ok {
		t.Errorf("Expected *MemoryTodoRepository, got %T", repo)
	}
}

func TestNewTodoRepository_UnsupportedType(t *testing.T) {
	// Create a config with an unsupported repository type
	cfg := config.DefaultConfig()
	cfg.Repository.Type = "unsupported"

	// Try to create the repository
	_, err := NewTodoRepository(cfg)
	if err == nil {
		t.Fatalf("Expected error for unsupported repository type, got nil")
	}
}

// Note: We're not testing the Supabase repository creation since it requires
// actual database connection details. This would be better tested in an
// integration test environment with a test database.
