package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/starbops/gottodo/pkg/config"
)

// SupabaseClient represents a client connection to the Supabase database
type SupabaseClient struct {
	DB *sql.DB
}

// NewSupabaseClient creates a new client connection to the Supabase database
func NewSupabaseClient(cfg *config.Config) (*SupabaseClient, error) {
	dbURL := cfg.GetSupabaseDBURL()
	if dbURL == "" {
		return nil, fmt.Errorf("Supabase database URL is not configured")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to Supabase PostgreSQL database")
	return &SupabaseClient{DB: db}, nil
}

// ExecWithLogging is a helper function that logs SQL queries and their parameters before execution
func (c *SupabaseClient) ExecWithLogging(query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Executing SQL: %s with args: %v", query, args)
	result, err := c.DB.Exec(query, args...)
	if err != nil {
		log.Printf("SQL Error: %v", err)
	}
	return result, err
}

// Close closes the database connection
func (c *SupabaseClient) Close() error {
	log.Println("Closing database connection")
	return c.DB.Close()
}

// ConnectToSupabase establishes a connection to the Supabase PostgreSQL database
// This is maintained for backward compatibility
func ConnectToSupabase(cfg *config.Config) (*sql.DB, error) {
	client, err := NewSupabaseClient(cfg)
	if err != nil {
		return nil, err
	}

	return client.DB, nil
}
