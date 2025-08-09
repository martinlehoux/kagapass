package models

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagapass/internal/types"
)

func TestFileSelectModelNew(t *testing.T) {
	dbList := types.DatabaseList{
		Databases: []types.Database{},
		LastUsed:  "",
	}
	
	model := NewFileSelectModel(dbList)
	if model == nil {
		t.Fatal("NewFileSelectModel returned nil")
	}
	
	if model.cursor != 0 {
		t.Error("Expected initial cursor position 0")
	}
	
	if model.inputMode {
		t.Error("Expected input mode to be false initially")
	}
}

func TestFileSelectModelWithDatabases(t *testing.T) {
	dbList := types.DatabaseList{
		Databases: []types.Database{
			{Name: "test1.kdbx", Path: "/path/to/test1.kdbx"},
			{Name: "test2.kdbx", Path: "/path/to/test2.kdbx"},
		},
		LastUsed: "/path/to/test1.kdbx",
	}
	
	model := NewFileSelectModel(dbList)
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
	
	model := NewFileSelectModel(dbList)
	
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
	model := NewFileSelectModel(types.DatabaseList{})
	
	// Enter input mode
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if !model.inputMode {
		t.Error("Expected input mode to be true after pressing 'a'")
	}
	
	if model.statusMessage == "" {
		t.Error("Expected status message when entering input mode")
	}
	
	// Type some text
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	
	if model.inputText != "test" {
		t.Errorf("Expected input text 'test', got '%s'", model.inputText)
	}
	
	// Test backspace
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if model.inputText != "tes" {
		t.Errorf("Expected input text 'tes' after backspace, got '%s'", model.inputText)
	}
	
	// Exit input mode with Esc
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.inputMode {
		t.Error("Expected input mode to be false after Esc")
	}
	
	if model.inputText != "" {
		t.Errorf("Expected empty input text after Esc, got '%s'", model.inputText)
	}
}

func TestFileSelectModelView(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{})
	
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
	
	model = NewFileSelectModel(dbList)
	view = model.View()
	
	if !strings.Contains(view, "test.kdbx") {
		t.Error("Expected view to contain database name")
	}
}

func TestSearchModelNew(t *testing.T) {
	model := NewSearchModel()
	if model == nil {
		t.Fatal("NewSearchModel returned nil")
	}
	
	if model.searchInput != "" {
		t.Error("Expected empty search input initially")
	}
	
	if model.cursor != 0 {
		t.Error("Expected cursor at 0 initially")
	}
}

func TestSearchModelSetEntries(t *testing.T) {
	model := NewSearchModel()
	
	entries := []types.Entry{
		{Title: "GitHub", Username: "user1"},
		{Title: "Gmail", Username: "user2"},
	}
	
	model.SetEntries(entries)
	
	if len(model.entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(model.entries))
	}
}

func TestSearchModelSearch(t *testing.T) {
	model := NewSearchModel()
	
	entries := []types.Entry{
		{Title: "GitHub Personal", Username: "user1"},
		{Title: "Gmail", Username: "user2"},
		{Title: "GitHub Work", Username: "user3"},
	}
	
	model.SetEntries(entries)
	
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
	model := NewSearchModel()
	
	entries := []types.Entry{
		{Title: "Entry1", Username: "user1"},
		{Title: "Entry2", Username: "user2"},
		{Title: "Entry3", Username: "user3"},
	}
	
	model.SetEntries(entries)
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

func TestDetailsModelNew(t *testing.T) {
	model := NewDetailsModel()
	if model == nil {
		t.Fatal("NewDetailsModel returned nil")
	}
	
	if model.entry != nil {
		t.Error("Expected no entry initially")
	}
}

func TestDetailsModelSetEntry(t *testing.T) {
	model := NewDetailsModel()
	
	entry := types.Entry{
		Title:    "Test Entry",
		Username: "testuser",
		Password: "testpass",
		URL:      "https://example.com",
		Notes:    "Test notes",
		Group:    "Test/Group",
	}
	
	model.SetEntry(entry)
	
	if model.entry == nil {
		t.Fatal("Expected entry to be set")
	}
	
	if model.entry.Title != entry.Title {
		t.Errorf("Expected title '%s', got '%s'", entry.Title, model.entry.Title)
	}
}

func TestDetailsModelView(t *testing.T) {
	model := NewDetailsModel()
	
	// Test view without entry
	view := model.View()
	if !strings.Contains(view, "No entry selected") {
		t.Error("Expected 'No entry selected' message when no entry is set")
	}
	
	// Test view with entry
	entry := types.Entry{
		Title:    "Test Entry",
		Username: "testuser",
		Password: "testpass",
		URL:      "https://example.com",
		Notes:    "Test notes",
		Group:    "Test/Group",
	}
	
	model.SetEntry(entry)
	view = model.View()
	
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

func TestPasswordModelNew(t *testing.T) {
	model := NewPasswordModel()
	if model == nil {
		t.Fatal("NewPasswordModel returned nil")
	}
	
	if model.password != "" {
		t.Error("Expected empty password initially")
	}
	
	if model.attempts != 0 {
		t.Error("Expected 0 attempts initially")
	}
}

func TestPasswordModelSetDatabase(t *testing.T) {
	model := NewPasswordModel()
	
	db := &types.Database{
		Name: "test.kdbx",
		Path: "/path/to/test.kdbx",
	}
	
	model.SetDatabase(db)
	
	if model.database != db {
		t.Error("Expected database to be set")
	}
	
	if model.password != "" {
		t.Error("Expected password to be cleared when setting database")
	}
}

func TestPasswordModelInput(t *testing.T) {
	model := NewPasswordModel()
	
	// Type password
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	
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
	model := NewPasswordModel()
	
	view := model.View()
	if !strings.Contains(view, "Enter Master Password") {
		t.Error("Expected view to contain password prompt")
	}
	
	// Set database
	db := &types.Database{
		Name: "test.kdbx",
		Path: "/path/to/test.kdbx",
	}
	model.SetDatabase(db)
	
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