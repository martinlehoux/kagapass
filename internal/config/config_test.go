package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/martinlehoux/kagapass/internal/types"
)

func TestNew(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if manager == nil {
		t.Fatal("New() returned nil manager")
	}

	// Check if config directory was created
	expectedConfigDir := filepath.Join(tmpDir, ".config", "kagapass")
	if _, err := os.Stat(expectedConfigDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created at %s", expectedConfigDir)
	}
}

func TestLoadConfigDefault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	config, err := manager.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Check default values
	expectedDefault := types.DefaultConfig()
	if config.ClipboardClearSeconds != expectedDefault.ClipboardClearSeconds {
		t.Errorf("Expected ClipboardClearSeconds %d, got %d", 
			expectedDefault.ClipboardClearSeconds, config.ClipboardClearSeconds)
	}

	if config.SearchDebounceMs != expectedDefault.SearchDebounceMs {
		t.Errorf("Expected SearchDebounceMs %d, got %d", 
			expectedDefault.SearchDebounceMs, config.SearchDebounceMs)
	}

	if config.MaxSearchResults != expectedDefault.MaxSearchResults {
		t.Errorf("Expected MaxSearchResults %d, got %d", 
			expectedDefault.MaxSearchResults, config.MaxSearchResults)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Create custom config
	customConfig := types.Config{
		ClipboardClearSeconds: 60,
		SearchDebounceMs:      200,
		MaxSearchResults:      100,
		SessionTimeoutHours:   24,
		DefaultDatabasePath:   "/test/path",
	}

	// Save config
	if err := manager.SaveConfig(customConfig); err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Load config
	loadedConfig, err := manager.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify values
	if loadedConfig.ClipboardClearSeconds != customConfig.ClipboardClearSeconds {
		t.Errorf("Expected ClipboardClearSeconds %d, got %d", 
			customConfig.ClipboardClearSeconds, loadedConfig.ClipboardClearSeconds)
	}

	if loadedConfig.DefaultDatabasePath != customConfig.DefaultDatabasePath {
		t.Errorf("Expected DefaultDatabasePath %s, got %s", 
			customConfig.DefaultDatabasePath, loadedConfig.DefaultDatabasePath)
	}
}

func TestLoadDatabaseListEmpty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	dbList, err := manager.LoadDatabaseList()
	if err != nil {
		t.Fatalf("LoadDatabaseList() failed: %v", err)
	}

	if len(dbList.Databases) != 0 {
		t.Errorf("Expected empty database list, got %d databases", len(dbList.Databases))
	}

	if dbList.LastUsed != "" {
		t.Errorf("Expected empty LastUsed, got %s", dbList.LastUsed)
	}
}

func TestSaveAndLoadDatabaseList(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Create test database list
	testDB := types.Database{
		Name: "test.kdbx",
		Path: "/path/to/test.kdbx",
	}

	dbList := types.DatabaseList{
		Databases: []types.Database{testDB},
		LastUsed:  testDB.Path,
	}

	// Save database list
	if err := manager.SaveDatabaseList(dbList); err != nil {
		t.Fatalf("SaveDatabaseList() failed: %v", err)
	}

	// Load database list
	loadedList, err := manager.LoadDatabaseList()
	if err != nil {
		t.Fatalf("LoadDatabaseList() failed: %v", err)
	}

	// Verify values
	if len(loadedList.Databases) != 1 {
		t.Errorf("Expected 1 database, got %d", len(loadedList.Databases))
	}

	if loadedList.Databases[0].Name != testDB.Name {
		t.Errorf("Expected database name %s, got %s", 
			testDB.Name, loadedList.Databases[0].Name)
	}

	if loadedList.LastUsed != testDB.Path {
		t.Errorf("Expected LastUsed %s, got %s", 
			testDB.Path, loadedList.LastUsed)
	}
}

func TestCorruptedConfigFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kagapass-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	manager, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// Write corrupted JSON to config file
	configPath := filepath.Join(tmpDir, ".config", "kagapass", "config.json")
	if err := os.WriteFile(configPath, []byte("invalid json{"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted config: %v", err)
	}

	// Should return default config when file is corrupted
	config, err := manager.LoadConfig()
	if err == nil {
		t.Error("Expected error when loading corrupted config, got nil")
	}

	// Should still have some default values even on error
	expectedDefault := types.DefaultConfig()
	if config.ClipboardClearSeconds != expectedDefault.ClipboardClearSeconds {
		t.Errorf("Expected default ClipboardClearSeconds %d, got %d", 
			expectedDefault.ClipboardClearSeconds, config.ClipboardClearSeconds)
	}
}