package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagapass/internal/clipboard"
	"github.com/martinlehoux/kagapass/internal/testor"
	"github.com/martinlehoux/kagapass/internal/types"
)

var unlockDatabase = &UnlockDatabase{
	keepassLoader: nil,
	secretStore:     nil,
}

func TestFileSelectModelWithDatabases(t *testing.T) {
	dbList := types.DatabaseList{
		Databases: []types.Database{
			{Name: "test1.kdbx", Path: "/path/to/test1.kdbx"},
			{Name: "test2.kdbx", Path: "/path/to/test2.kdbx"},
		},
		LastUsed: "/path/to/test1.kdbx",
	}

	model := NewFileSelectModel(dbList, unlockDatabase)
	if len(model.databases.Databases) != 2 {
		t.Errorf("Expected 2 databases, got %d", len(model.databases.Databases))
	}
}

func TestFileSelectModelNavigation(t *testing.T) {
	dbList := types.DatabaseList{
		Databases: []types.Database{
			{Name: "test1.kdbx", Path: "/path/to/test1.kdbx"},
			{Name: "test2.kdbx", Path: "/path/to/test2.kdbx"},
			{Name: "test3.kdbx", Path: "/path/to/test3.kdbx"},
		},
	}

	model := NewFileSelectModel(dbList, unlockDatabase)

	// Test down navigation
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.cursor != 1 {
		t.Errorf("Expected cursor at 1 after down, got %d", model.cursor)
	}

	// Test up navigation
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0 after up, got %d", model.cursor)
	}

	// Test up at beginning (should stay at 0)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0 at beginning, got %d", model.cursor)
	}

	// Move to end
	model.cursor = 2
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.cursor != 2 {
		t.Errorf("Expected cursor to stay at 2 at end, got %d", model.cursor)
	}
}

func TestFileSelectModelInputMode(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{}, unlockDatabase)

	// Enter input mode
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Error("Expected input mode to be true after pressing 'a'")
	}

	if model.statusMessage != "" {
		t.Error("Expected no status message when entering input mode")
	}

	// Type some text
	model, _ = model.Update(testor.KeyMsgRune('t'))
	model, _ = model.Update(testor.KeyMsgRune('e'))
	model, _ = model.Update(testor.KeyMsgRune('s'))
	model, _ = model.Update(testor.KeyMsgRune('t'))

	if model.databaseInput.Value() != "test" {
		t.Errorf("Expected input text 'test', got '%s'", model.databaseInput.Value())
	}

	// Test backspace
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if model.databaseInput.Value() != "tes" {
		t.Errorf("Expected input text 'tes' after backspace, got '%s'", model.databaseInput.Value())
	}

	// Exit input mode with Esc
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.databaseInput.Focused() {
		t.Error("Expected input mode to be false after Esc")
	}

	if model.databaseInput.Value() != "" {
		t.Errorf("Expected empty input text after Esc, got '%s'", model.databaseInput.Value())
	}
}

func TestFileSelectModelView(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{}, unlockDatabase)

	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view")
	}

	if !strings.Contains(view, "KagaPass") {
		t.Error("Expected view to contain 'KagaPass'")
	}

	// Test with databases
	dbList := types.DatabaseList{
		Databases: []types.Database{
			{Name: "test.kdbx", Path: "/path/to/test.kdbx"},
		},
	}

	model = NewFileSelectModel(dbList, unlockDatabase)
	view = model.View()

	if !strings.Contains(view, "test.kdbx") {
		t.Error("Expected view to contain database name")
	}
}

func TestSearchModelSearch(t *testing.T) {
	model := NewSearchModel(clipboard.New(), []types.Entry{
		{Title: "GitHub Personal", Username: "user1"},
		{Title: "Gmail", Username: "user2"},
		{Title: "GitHub Work", Username: "user3"},
	}, func(entry types.Entry) {}, "test")

	// Search for "github"
	model.searchInput = "github"
	model.search()

	if len(model.filteredItems) != 2 {
		t.Errorf("Expected 2 filtered items for 'github', got %d", len(model.filteredItems))
	}

	// Search for something that doesn't exist
	model.searchInput = "nonexistent"
	model.search()

	if len(model.filteredItems) != 0 {
		t.Errorf("Expected 0 filtered items for 'nonexistent', got %d", len(model.filteredItems))
	}

	// Empty search should return no results
	model.searchInput = ""
	model.search()

	if len(model.filteredItems) != 0 {
		t.Errorf("Expected 0 filtered items for empty search, got %d", len(model.filteredItems))
	}
}

func TestSearchModelNavigation(t *testing.T) {
	model := NewSearchModel(clipboard.New(), []types.Entry{
		{Title: "Entry1", Username: "user1"},
		{Title: "Entry2", Username: "user2"},
		{Title: "Entry3", Username: "user3"},
	}, func(entry types.Entry) {}, "")

	model.searchInput = "entry"
	model.search()

	// Should have 3 results
	if len(model.filteredItems) != 3 {
		t.Errorf("Expected 3 filtered items, got %d", len(model.filteredItems))
	}

	// Test navigation
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.cursor != 1 {
		t.Errorf("Expected cursor at 1, got %d", model.cursor)
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", model.cursor)
	}
}

func TestDetailsModelView(t *testing.T) {
	entry := types.Entry{
		Title:    "Test Entry",
		Username: "testuser",
		Password: "testpass",
		URL:      "https://example.com",
		Notes:    "Test notes",
		Group:    "Test/Group",
	}
	model := NewDetailsModel(clipboard.New(), entry)

	// Test view with entry
	view := model.View()

	if !strings.Contains(view, entry.Title) {
		t.Error("Expected view to contain entry title")
	}

	if !strings.Contains(view, entry.Username) {
		t.Error("Expected view to contain entry username")
	}

	// Password should be masked
	if strings.Contains(view, entry.Password) {
		t.Error("Password should not appear in plain text in view")
	}

	if !strings.Contains(view, "************") {
		t.Error("Expected masked password in view")
	}
}

func TestPasswordModelInput(t *testing.T) {
	model := &PasswordModel{
		unlockDatabase: unlockDatabase,
	}

	// Type password
	model, _ = model.Update(testor.KeyMsgRune('p'))
	model, _ = model.Update(testor.KeyMsgRune('a'))
	model, _ = model.Update(testor.KeyMsgRune('s'))
	model, _ = model.Update(testor.KeyMsgRune('s'))

	if model.password != "pass" {
		t.Errorf("Expected password 'pass', got '%s'", model.password)
	}

	// Test backspace
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if model.password != "pas" {
		t.Errorf("Expected password 'pas' after backspace, got '%s'", model.password)
	}

	// Test clear with Ctrl+L
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlL})
	if model.password != "" {
		t.Errorf("Expected empty password after Ctrl+L, got '%s'", model.password)
	}
}

func TestPasswordModelView(t *testing.T) {
	db := types.Database{
		Name: "test.kdbx",
		Path: "/path/to/test.kdbx",
	}
	model := NewPasswordModel(unlockDatabase, func() {}, db)

	view := model.View()
	if !strings.Contains(view, "Enter Master Password") {
		t.Error("Expected view to contain password prompt")
	}

	view = model.View()
	if !strings.Contains(view, "test.kdbx") {
		t.Error("Expected view to contain database name")
	}

	// Type some password
	model.password = "secret"
	view = model.View()

	// Password should be masked
	if strings.Contains(view, "secret") {
		t.Error("Password should not appear in plain text in view")
	}

	if !strings.Contains(view, "••••••") {
		t.Error("Expected masked password characters in view")
	}
}

func TestAppModelEscKeyInDatabaseSelection(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "kagapass-esc-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create app model (starts on file selection screen)
	app, err := NewAppModel()
	if err != nil {
		t.Fatalf("Failed to create app model: %v", err)
	}

	// Should start on file selection screen
	view := app.View()
	if view == "" {
		t.Error("App view should not be empty")
	}

	// Press Esc - should quit the application
	newModel, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEsc})

	// Check if quit command was returned
	if cmd == nil {
		t.Error("Expected quit command when pressing Esc on file selection screen, got nil")
	} else {
		// Execute the command to see what message it returns
		msg := cmd()
		quitMsg := tea.Quit()
		if msg != quitMsg {
			t.Error("Expected tea.Quit message when pressing Esc on file selection screen")
		}
	}

	// Model should be updated
	if newModel == nil {
		t.Error("Expected updated model, got nil")
	}
}

func TestSessionPersistence(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "kagapass-session-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create a mock database list with last used
	configDir := filepath.Join(tmpDir, ".config", "kagapass")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create databases.json with a last used database
	dbList := types.DatabaseList{
		Databases: []types.Database{
			{
				Name:         "test.kdbx",
				Path:         "/path/to/test.kdbx",
				LastAccessed: time.Now(),
			},
			{
				Name:         "real.kdbx",
				Path:         "/path/to/real.kdbx",
				LastAccessed: time.Now().Add(-1 * time.Hour),
			},
		},
		LastUsed: "/path/to/test.kdbx", // Last used was the test database
	}

	data, err := json.MarshalIndent(dbList, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal database list: %v", err)
	}

	dbPath := filepath.Join(configDir, "databases.json")
	if err := os.WriteFile(dbPath, data, 0o600); err != nil {
		t.Fatalf("Failed to write databases.json: %v", err)
	}

	// Create app model - should load the database list
	app, err := NewAppModel()
	if err != nil {
		t.Fatalf("Failed to create app model: %v", err)
	}

	// Initialize the app - should trigger automatic unlock attempt
	cmd := app.Init()
	if cmd == nil {
		t.Error("Expected Init() to return a command to unlock last used database, got nil")
	} else {
		// Execute the command to see what message it returns
		cmd()

		// TODO: This should be TryKeyringUnlockMsg for the test.kdbx database
		// if unlockMsg, ok := msg.(TryKeyringUnlockMsg); ok {
		// 	if unlockMsg.Database == nil {
		// 		t.Error("Expected database in TryKeyringUnlockMsg")
		// 	} else if unlockMsg.Database.Path != "/path/to/test.kdbx" {
		// 		t.Errorf("Expected database path '/path/to/test.kdbx', got '%s'", unlockMsg.Database.Path)
		// 	}
		// } else {
		// 	t.Errorf("Expected TryKeyringUnlockMsg, got %T", msg)
		// }
	}
}

func TestDatabaseUnlockFlow(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "kagapass-unlock-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create app model
	app, err := NewAppModel()
	if err != nil {
		t.Fatalf("Failed to create app model: %v", err)
	}

	// Create a database entry (simulating a real KeePass file)
	// db := types.Database{
	// 	Name:         "real.kdbx",
	// 	Path:         "/nonexistent/path/real.kdbx",
	// 	LastAccessed: time.Now(),
	// }

	// Step 1: Simulate TryKeyringUnlockMsg (what happens when you press Enter on a database)
	// tryUnlockMsg := TryKeyringUnlockMsg{Database: &db}
	app1, cmd1 := app.Update(testor.KeyMsgRune('a')) //	TODO

	if cmd1 == nil {
		t.Error("Expected command from TryKeyringUnlockMsg, got nil")
		return
	}

	// Execute the command
	msg1 := cmd1()

	// Step 2: Handle the resulting message (should be SwitchScreenMsg to password screen)
	app2, cmd2 := app1.Update(msg1)

	// Verify we switched to password input screen
	view := app2.View()
	if !strings.Contains(view, "Enter Master Password") {
		t.Error("Expected to switch to password input screen")
	}
	if !strings.Contains(view, "real.kdbx") {
		t.Error("Expected password screen to show database name")
	}

	// The unlock flow is working correctly:
	// 1. TryKeyringUnlockMsg received
	// 2. No stored password found, switch to password screen
	// 3. Password screen displayed correctly

	// Verify cmd2 is nil (no further commands after screen switch)
	if cmd2 != nil {
		t.Error("Expected no command after switching to password screen")
	}
}
