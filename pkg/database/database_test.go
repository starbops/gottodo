package database

import (
	"testing"

	"github.com/starbops/gottodo/pkg/config"
)

func TestNewSupabaseClient_NoDBURL(t *testing.T) {
	// Create a config with empty database URL
	cfg := config.DefaultConfig()
	cfg.Database.SupabaseDBURL = ""

	// Try to create a client without the database URL
	client, err := NewSupabaseClient(cfg)

	// Should get an error
	if err == nil {
		t.Error("Expected error when database URL is not set, got nil")
		if client != nil && client.DB != nil {
			client.DB.Close()
		}
	}

	// Error message should mention the missing database URL
	if err != nil && err.Error() != "Supabase database URL is not configured" {
		t.Errorf("Expected error about missing database URL, got: %v", err)
	}
}

func TestConnectToSupabase_NoDBURL(t *testing.T) {
	// Create a config with empty database URL
	cfg := config.DefaultConfig()
	cfg.Database.SupabaseDBURL = ""

	// Try to connect without the database URL
	db, err := ConnectToSupabase(cfg)

	// Should get an error
	if err == nil {
		t.Error("Expected error when database URL is not set, got nil")
		if db != nil {
			db.Close()
		}
	}

	// Error message should mention the missing database URL
	if err != nil && err.Error() != "Supabase database URL is not configured" {
		t.Errorf("Expected error about missing database URL, got: %v", err)
	}
}

// Note: We're not testing actual database connections here since that would require
// a real database. In a more comprehensive test suite, you might want to use
// a test database or a mock.
