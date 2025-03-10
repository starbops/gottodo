package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Repository.Type != MemoryRepository {
		t.Errorf("Expected default repository type to be %s, got %s", MemoryRepository, cfg.Repository.Type)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port to be 8080, got %s", cfg.Server.Port)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Test creating a new config file
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Repository.Type != MemoryRepository {
		t.Errorf("Expected default repository type to be %s, got %s", MemoryRepository, cfg.Repository.Type)
	}

	// Verify the file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created")
	}

	// Modify and save the config
	cfg.Repository.Type = SupabaseRepository
	cfg.Server.Port = "9000"

	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the modified config
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load modified config: %v", err)
	}

	if loadedCfg.Repository.Type != SupabaseRepository {
		t.Errorf("Expected repository type to be %s, got %s", SupabaseRepository, loadedCfg.Repository.Type)
	}

	if loadedCfg.Server.Port != "9000" {
		t.Errorf("Expected port to be 9000, got %s", loadedCfg.Server.Port)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")

	// Write invalid JSON to the file
	if err := os.WriteFile(configPath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Try to load the invalid config
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Fatalf("Expected error when loading invalid JSON, got nil")
	}
}
