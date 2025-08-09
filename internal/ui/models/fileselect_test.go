package models

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagapass/internal/testor"
	"github.com/martinlehoux/kagapass/internal/types"
)

func TestInputModeNavigationKeyConflicts(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{}, func(database types.Database, password string) tea.Cmd { return nil })

	// Enter input mode by pressing 'a'
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Fatal("Should be in input mode after pressing 'a'")
	}

	// Test typing letters that might conflict with navigation
	testLetters := []rune{'k', 'a', 'g', 'j', 'h', 'l'}

	for _, letter := range testLetters {
		initialLength := len(model.databaseInput.Value())
		model, _ = model.Update(testor.KeyMsgRune(letter))

		// Check if letter was added to input
		if len(model.databaseInput.Value()) != initialLength+1 {
			t.Errorf("Letter '%c' was not added to input. Expected length %d, got %d. Input: '%s'",
				letter, initialLength+1, len(model.databaseInput.Value()), model.databaseInput.Value())
		}

		// Verify the letter is actually the one we typed
		value := model.databaseInput.Value()
		if len(value) > 0 && value[len(value)-1] != byte(letter) {
			t.Errorf("Last character should be '%c', but got '%c'",
				letter, value[len(value)-1])
		}

		// Make sure we're still in input mode
		if !model.databaseInput.Focused() {
			t.Errorf("Should still be in input mode after typing '%c'", letter)
		}
	}

	// Final check - input should contain all the letters except the first 'a' which was used to enter input mode
	expectedInput := "kagjhl"
	value := model.databaseInput.Value()
	if value != expectedInput {
		t.Errorf("Expected input text '%s', got '%s'", expectedInput, value)
	}
}

func TestNavigationKeysDisabledInInputMode(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{
		Databases: []types.Database{
			{Name: "test1.kdbx", Path: "/path1"},
			{Name: "test2.kdbx", Path: "/path2"},
		},
	}, func(database types.Database, password string) tea.Cmd { return nil })

	// Enter input mode
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Fatal("Should be in input mode")
	}

	initialCursor := model.cursor

	// Test navigation keys that should NOT work in input mode
	navigationKeys := []tea.KeyType{
		tea.KeyUp,
		tea.KeyDown,
	}

	for _, key := range navigationKeys {
		model, _ = model.Update(tea.KeyMsg{Type: key})

		// Cursor should not change when in input mode
		if model.cursor != initialCursor {
			t.Errorf("Cursor should not change in input mode for key %v. Was %d, now %d",
				key, initialCursor, model.cursor)
		}

		// Should still be in input mode
		if !model.databaseInput.Focused() {
			t.Errorf("Should still be in input mode after pressing %v", key)
		}
	}

	// Test vim-style navigation keys that should be treated as regular input
	vimKeys := []rune{'j', 'k', 'h', 'l'}

	for _, key := range vimKeys {
		initialLength := len(model.databaseInput.Value())
		initialCursor := model.cursor

		model, _ = model.Update(testor.KeyMsgRune(key))

		// Should add to input text
		if len(model.databaseInput.Value()) != initialLength+1 {
			t.Errorf("Vim key '%c' should be added as input, not used for navigation", key)
		}

		// Cursor position should not change for navigation
		if model.cursor != initialCursor {
			t.Errorf("Cursor should not move when typing '%c' in input mode", key)
		}
	}
}

func TestInputModeEscapeBehavior(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{}, func(database types.Database, password string) tea.Cmd { return nil })

	// Enter input mode
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Fatal("Should be in input mode")
	}

	// Type some text including problematic letters
	testText := "kagapass_database"
	for _, char := range testText {
		model, _ = model.Update(testor.KeyMsgRune(char))
	}

	if model.databaseInput.Value() != testText {
		t.Errorf("Expected input text '%s', got '%s'", testText, model.databaseInput.Value())
	}

	// Exit input mode with Esc
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if model.databaseInput.Focused() {
		t.Error("Should exit input mode with Esc")
	}

	if model.databaseInput.Value() != "" {
		t.Errorf("Input text should be cleared on Esc, got '%s'", model.databaseInput.Value())
	}
}

func TestNavigationWorksWhenNotInInputMode(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{
		Databases: []types.Database{
			{Name: "test1.kdbx", Path: "/path1"},
			{Name: "test2.kdbx", Path: "/path2"},
			{Name: "test3.kdbx", Path: "/path3"},
		},
	}, func(database types.Database, password string) tea.Cmd { return nil })

	// Test vim-style navigation works in normal mode
	initialCursor := model.cursor

	// Test 'j' for down
	model, _ = model.Update(testor.KeyMsgRune('j'))
	if model.cursor != initialCursor+1 {
		t.Errorf("'j' should move cursor down. Expected %d, got %d", initialCursor+1, model.cursor)
	}

	// Test 'k' for up
	model, _ = model.Update(testor.KeyMsgRune('k'))
	if model.cursor != initialCursor {
		t.Errorf("'k' should move cursor up. Expected %d, got %d", initialCursor, model.cursor)
	}

	// Make sure we're not in input mode
	if model.databaseInput.Focused() {
		t.Error("Should not be in input mode during navigation")
	}
}

func TestFileSelectInputModeToggling(t *testing.T) {
	model := NewFileSelectModel(types.DatabaseList{}, func(database types.Database, password string) tea.Cmd { return nil })

	// Initially not in input mode
	if model.databaseInput.Focused() {
		t.Error("Should not be in input mode initially")
	}

	// Enter input mode with 'a'
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Error("Should be in input mode after pressing 'a'")
	}

	// Exit with Esc
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.databaseInput.Focused() {
		t.Error("Should exit input mode with Esc")
	}

	// Enter again
	model, _ = model.Update(testor.KeyMsgRune('a'))
	if !model.databaseInput.Focused() {
		t.Error("Should be able to re-enter input mode")
	}
}
