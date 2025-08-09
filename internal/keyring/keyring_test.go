package keyring

import (
	"testing"
)

func TestStoreAndRetrieve(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	testPath := "/test/path/database.kdbx"
	testPassword := "test_password_123"

	// Store password
	err = manager.Store(testPath, testPassword)
	if err != nil {
		t.Skipf("Failed to store in keyring (expected in test env): %v", err)
	}

	// Retrieve password
	retrieved, err := manager.Retrieve(testPath)
	if err != nil {
		t.Errorf("Failed to retrieve password: %v", err)
	}

	if retrieved != testPassword {
		t.Errorf("Expected password '%s', got '%s'", testPassword, retrieved)
	}

	// Clean up
	manager.Remove(testPath)
}

func TestRetrieveNonexistent(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	// Try to retrieve non-existent password
	_, err = manager.Retrieve("/nonexistent/path.kdbx")
	if err == nil {
		t.Error("Expected error when retrieving non-existent password")
	}
}

func TestExists(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	testPath := "/test/exists/database.kdbx"

	// Should not exist initially
	if manager.Exists(testPath) {
		t.Error("Password should not exist initially")
	}

	// Store password
	err = manager.Store(testPath, "test_password")
	if err != nil {
		t.Skipf("Failed to store in keyring: %v", err)
	}

	// Should exist now
	if !manager.Exists(testPath) {
		t.Error("Password should exist after storing")
	}

	// Clean up
	manager.Remove(testPath)

	// Should not exist after removal
	if manager.Exists(testPath) {
		t.Error("Password should not exist after removal")
	}
}

func TestRemoveNonexistent(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	// Try to remove non-existent password (should not crash)
	err = manager.Remove("/nonexistent/path.kdbx")
	if err == nil {
		t.Error("Expected error when removing non-existent password")
	}
}

func TestListStored(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	// Initially should be empty or contain other entries
	initialPaths, err := manager.ListStored()
	if err != nil {
		t.Skipf("Failed to list stored passwords: %v", err)
	}

	testPath := "/test/list/database.kdbx"

	// Store a test password
	err = manager.Store(testPath, "test_password")
	if err != nil {
		t.Skipf("Failed to store password: %v", err)
	}

	// List should now contain our test path
	paths, err := manager.ListStored()
	if err != nil {
		t.Errorf("Failed to list stored passwords: %v", err)
	}

	found := false
	for _, path := range paths {
		if path == testPath {
			found = true
			break
		}
	}

	if !found {
		t.Error("Stored password path not found in list")
	}

	// Should have more entries than initially
	if len(paths) <= len(initialPaths) {
		t.Error("Expected more entries after storing password")
	}

	// Clean up
	manager.Remove(testPath)
}

func TestStoreEmptyPassword(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	testPath := "/test/empty/database.kdbx"

	// Store empty password
	err = manager.Store(testPath, "")
	if err != nil {
		t.Skipf("Failed to store empty password: %v", err)
	}

	// Retrieve empty password
	retrieved, err := manager.Retrieve(testPath)
	if err != nil {
		t.Errorf("Failed to retrieve empty password: %v", err)
	}

	if retrieved != "" {
		t.Errorf("Expected empty password, got '%s'", retrieved)
	}

	// Clean up
	manager.Remove(testPath)
}

func TestStorePathWithSpecialChars(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	// Test path with special characters
	testPath := "/test/special chars & symbols/data base.kdbx"
	testPassword := "password_with_special_chars!@#$%^&*()"

	err = manager.Store(testPath, testPassword)
	if err != nil {
		t.Skipf("Failed to store password with special chars: %v", err)
	}

	retrieved, err := manager.Retrieve(testPath)
	if err != nil {
		t.Errorf("Failed to retrieve password with special chars: %v", err)
	}

	if retrieved != testPassword {
		t.Errorf("Expected password '%s', got '%s'", testPassword, retrieved)
	}

	// Clean up
	manager.Remove(testPath)
}

func TestMultipleOperations(t *testing.T) {
	manager, err := New()
	if err != nil {
		t.Skipf("Keyring not available in test environment: %v", err)
	}

	// Store multiple passwords
	paths := []string{
		"/test/multi/db1.kdbx",
		"/test/multi/db2.kdbx",
		"/test/multi/db3.kdbx",
	}

	passwords := []string{
		"password1",
		"password2",
		"password3",
	}

	// Store all
	for i, path := range paths {
		err = manager.Store(path, passwords[i])
		if err != nil {
			t.Skipf("Failed to store password %d: %v", i, err)
		}
	}

	// Verify all exist
	for i, path := range paths {
		if !manager.Exists(path) {
			t.Errorf("Password %d should exist", i)
		}
	}

	// Retrieve all
	for i, path := range paths {
		retrieved, err := manager.Retrieve(path)
		if err != nil {
			t.Errorf("Failed to retrieve password %d: %v", i, err)
		}
		if retrieved != passwords[i] {
			t.Errorf("Password %d mismatch: expected '%s', got '%s'",
				i, passwords[i], retrieved)
		}
	}

	// Clean up all
	for _, path := range paths {
		manager.Remove(path)
	}
}
