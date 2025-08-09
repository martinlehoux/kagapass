package database

import (
	"os"
	"testing"
	"time"

	"github.com/martinlehoux/kagapass/internal/types"
)

func TestNew(t *testing.T) {
	manager := New()
	if manager == nil {
		t.Fatal("New() returned nil")
	}

	if manager.IsOpen() {
		t.Error("New manager should not be open")
	}

	if manager.GetFilePath() != "" {
		t.Error("New manager should have empty file path")
	}
}

func TestOpenNonexistentFile(t *testing.T) {
	manager := New()
	
	err := manager.Open("/nonexistent/file.kdbx", "password")
	if err == nil {
		t.Error("Expected error when opening nonexistent file, got nil")
	}

	if manager.IsOpen() {
		t.Error("Manager should not be open after failed open")
	}
}

func TestOpenWithEmptyPath(t *testing.T) {
	manager := New()
	
	err := manager.Open("", "password")
	if err == nil {
		t.Error("Expected error when opening empty path, got nil")
	}
}

func TestOpenWithEmptyPassword(t *testing.T) {
	// Create a temporary file (not a real KeePass file)
	tmpFile, err := os.CreateTemp("", "test-*.kdbx")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	manager := New()
	
	err = manager.Open(tmpFile.Name(), "")
	if err == nil {
		t.Error("Expected error when opening with empty password, got nil")
	}
}

func TestGetEntriesWithoutOpen(t *testing.T) {
	manager := New()
	
	entries, err := manager.GetEntries()
	if err == nil {
		t.Error("Expected error when getting entries without opening database, got nil")
	}

	if entries != nil {
		t.Error("Expected nil entries when database not open")
	}
}

func TestCloseUnopened(t *testing.T) {
	manager := New()
	
	// Should not panic
	manager.Close()
	
	if manager.IsOpen() {
		t.Error("Manager should not be open after close")
	}
}

func TestGetFileName(t *testing.T) {
	manager := New()
	
	// Test empty path
	if filename := manager.GetFileName(); filename != "" {
		t.Errorf("Expected empty filename, got %s", filename)
	}

	// Test with simulated path (without actually opening)
	manager.filePath = "/path/to/database.kdbx"
	if filename := manager.GetFileName(); filename != "database.kdbx" {
		t.Errorf("Expected 'database.kdbx', got %s", filename)
	}
}

func TestCreateTestEntries(t *testing.T) {
	entries := CreateTestEntries()
	
	if len(entries) == 0 {
		t.Error("CreateTestEntries() returned empty slice")
	}

	// Check that entries have required fields
	for i, entry := range entries {
		if entry.Title == "" {
			t.Errorf("Entry %d has empty title", i)
		}
		
		if entry.Username == "" {
			t.Errorf("Entry %d has empty username", i)
		}
		
		if entry.Password == "" {
			t.Errorf("Entry %d has empty password", i)
		}
		
		// Check timestamps are reasonable
		now := time.Now()
		if entry.Created.After(now) {
			t.Errorf("Entry %d has future creation time", i)
		}
		
		if entry.Modified.After(now) {
			t.Errorf("Entry %d has future modification time", i)
		}
		
		if entry.Modified.Before(entry.Created) {
			t.Errorf("Entry %d has modification time before creation time", i)
		}
	}
}

func TestGetTestEntries(t *testing.T) {
	entries1 := CreateTestEntries()
	entries2 := GetTestEntries()
	
	if len(entries1) != len(entries2) {
		t.Error("CreateTestEntries() and GetTestEntries() return different lengths")
	}
	
	// Should be the same data
	for i := range entries1 {
		if entries1[i].Title != entries2[i].Title {
			t.Errorf("Entry %d titles differ: %s vs %s", 
				i, entries1[i].Title, entries2[i].Title)
		}
	}
}

func TestCollectEntriesFromGroupNil(t *testing.T) {
	manager := New()
	var entries []types.Entry
	
	// Should not panic with nil group
	manager.collectEntriesFromGroup(nil, "", &entries)
	
	if len(entries) != 0 {
		t.Error("Expected no entries from nil group")
	}
}

func TestEntryGroupPaths(t *testing.T) {
	entries := CreateTestEntries()
	
	// Check that group paths are properly formatted
	for i, entry := range entries {
		if entry.Group == "" {
			t.Errorf("Entry %d has empty group path", i)
		}
		
		// Group paths should not start or end with "/"
		if len(entry.Group) > 0 {
			if entry.Group[0] == '/' {
				t.Errorf("Entry %d group path starts with '/': %s", i, entry.Group)
			}
			if entry.Group[len(entry.Group)-1] == '/' {
				t.Errorf("Entry %d group path ends with '/': %s", i, entry.Group)
			}
		}
	}
}

func TestDuplicateTestEntries(t *testing.T) {
	entries := CreateTestEntries()
	
	// Check for duplicate titles (should not have any)
	titleMap := make(map[string]int)
	for i, entry := range entries {
		if prevIndex, exists := titleMap[entry.Title]; exists {
			t.Errorf("Duplicate title '%s' found at indices %d and %d", 
				entry.Title, prevIndex, i)
		}
		titleMap[entry.Title] = i
	}
}

func TestEntryFieldValidation(t *testing.T) {
	entries := CreateTestEntries()
	
	for i, entry := range entries {
		// URL should be valid format if present
		if entry.URL != "" {
			if !isValidURLFormat(entry.URL) {
				t.Errorf("Entry %d has invalid URL format: %s", i, entry.URL)
			}
		}
		
		// Username should not contain obvious invalid characters
		if containsInvalidChars(entry.Username) {
			t.Errorf("Entry %d has suspicious username: %s", i, entry.Username)
		}
	}
}

// Helper functions for validation
func isValidURLFormat(url string) bool {
	// Simple validation - should start with http:// or https://
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

func containsInvalidChars(username string) bool {
	// Check for obviously invalid characters in usernames
	invalidChars := []string{"\n", "\r", "\t", "\x00"}
	for _, char := range invalidChars {
		if len(char) > 0 && len(username) > 0 && username[0] == char[0] {
			return true
		}
	}
	return false
}