package database

import (
	"os"
	"testing"
)

func TestNewSupabaseClient_NoEnvVar(t *testing.T) {
	// Save current env var value and restore it after the test
	oldDBURL := os.Getenv("SUPABASE_DB_URL")
	defer os.Setenv("SUPABASE_DB_URL", oldDBURL)

	// Clear the env var
	os.Unsetenv("SUPABASE_DB_URL")

	// Try to create a client without the env var
	client, err := NewSupabaseClient()

	// Should get an error
	if err == nil {
		t.Error("Expected error when SUPABASE_DB_URL is not set, got nil")
		if client != nil && client.DB != nil {
			client.DB.Close()
		}
	}

	// Error message should mention the missing env var
	if err != nil && err.Error() != "SUPABASE_DB_URL environment variable is not set" {
		t.Errorf("Expected error about missing env var, got: %v", err)
	}
}

func TestConnectToSupabase_NoEnvVar(t *testing.T) {
	// Save current env var value and restore it after the test
	oldDBURL := os.Getenv("SUPABASE_DB_URL")
	defer os.Setenv("SUPABASE_DB_URL", oldDBURL)

	// Clear the env var
	os.Unsetenv("SUPABASE_DB_URL")

	// Try to connect without the env var
	db, err := ConnectToSupabase()

	// Should get an error
	if err == nil {
		t.Error("Expected error when SUPABASE_DB_URL is not set, got nil")
		if db != nil {
			db.Close()
		}
	}

	// Error message should mention the missing env var
	if err != nil && err.Error() != "SUPABASE_DB_URL environment variable is not set" {
		t.Errorf("Expected error about missing env var, got: %v", err)
	}
}

// Note: We're not testing actual database connections here since that would require
// a real database. In a more comprehensive test suite, you might want to use
// a test database or a mock.
