package repositories

import (
	"fmt"
	"log"

	"github.com/starbops/gottodo/pkg/config"
	"github.com/starbops/gottodo/pkg/database"
)

// NewTodoRepository creates a TodoRepository based on the provided configuration
func NewTodoRepository(cfg *config.Config) (TodoRepository, error) {
	switch cfg.Repository.Type {
	case config.MemoryRepository:
		log.Println("Using in-memory todo repository")
		return NewMemoryTodoRepository(), nil

	case config.SupabaseRepository:
		log.Println("Using Supabase todo repository")
		// Connect to Supabase
		db, err := database.ConnectToSupabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Supabase: %w", err)
		}

		return NewSupabaseTodoRepository(db), nil

	default:
		return nil, fmt.Errorf("unsupported repository type: %s", cfg.Repository.Type)
	}
}
