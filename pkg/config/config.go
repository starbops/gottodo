package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RepositoryType defines the type of repository to use
type RepositoryType string

const (
	// MemoryRepository uses in-memory storage (for development/testing)
	MemoryRepository RepositoryType = "memory"

	// SupabaseRepository uses Supabase as storage (for production)
	SupabaseRepository RepositoryType = "supabase"
)

// Config represents the application configuration
type Config struct {
	// Repository configuration
	Repository struct {
		// Type is the repository type to use (memory or supabase)
		Type RepositoryType `json:"type"`
	} `json:"repository"`

	// Server configuration
	Server struct {
		// Port is the port number the server will listen on
		Port string `json:"port"`
	} `json:"server"`

	// Additional configuration settings can be added here in the future
	// For example:
	// Auth struct {
	//     SessionTimeout int `json:"sessionTimeout"` // in minutes
	// } `json:"auth"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{}

	// Set default repository type to memory
	cfg.Repository.Type = MemoryRepository

	// Set default server port
	cfg.Server.Port = "8080"

	return cfg
}

// LoadConfig loads the configuration from the specified file
// If the file doesn't exist, it returns the default configuration
func LoadConfig(configPath string) (*Config, error) {
	// Start with default config
	cfg := DefaultConfig()

	// Try to read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, create it with default values
		if os.IsNotExist(err) {
			// Ensure directory exists
			dir := filepath.Dir(configPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}

			// Write default config to file
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return nil, err
			}

			if err := os.WriteFile(configPath, data, 0644); err != nil {
				return nil, err
			}

			return cfg, nil
		}

		return nil, err
	}

	// Parse the config file
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the specified file
func SaveConfig(cfg *Config, configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal the config to JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write the config to file
	return os.WriteFile(configPath, data, 0644)
}
