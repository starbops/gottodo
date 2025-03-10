package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectToSupabase establishes a connection to the Supabase PostgreSQL database
func ConnectToSupabase() (*sql.DB, error) {
	dbURL := os.Getenv("SUPABASE_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("SUPABASE_DB_URL environment variable is not set")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60) // 5 minutes

	log.Println("Connected to Supabase PostgreSQL database")
	return db, nil
}
