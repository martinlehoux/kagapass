package clipboard

import (
	"testing"
	"time"
)

func TestCopyWithoutClearTime(t *testing.T) {
	manager := New()

	// Test copying without auto-clear (0 duration)
	err := manager.Copy("test content", 0)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	// Should not have a timer set
	if manager.clearTimer != nil {
		t.Error("Expected no timer when clearAfter is 0")
	}
}

func TestCopyWithClearTime(t *testing.T) {
	manager := New()
	testContent := "test content for auto-clear"

	// Test copying with auto-clear
	err := manager.Copy(testContent, 100*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	// Should have a timer set
	if manager.clearTimer == nil {
		t.Error("Expected timer when clearAfter > 0")
	}

	// Verify content is there initially
	initial, err := manager.Get()
	if err != nil {
		t.Skipf("Can't read clipboard to verify test: %v", err)
	}

	if initial != testContent {
		t.Errorf("Expected initial content '%s', got '%s'", testContent, initial)
	}

	// Wait for timer to complete
	time.Sleep(150 * time.Millisecond)

	// Verify clipboard was actually cleared
	after, err := manager.Get()
	if err != nil {
		t.Errorf("Failed to read clipboard after clear: %v", err)
	} else if after == testContent {
		t.Error("Clipboard should have been cleared after timeout, but content is still there")
	}
}

func TestCopyOverwrite(t *testing.T) {
	manager := New()

	// First copy with timer
	err := manager.Copy("first content", 500*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	firstTimer := manager.clearTimer
	if firstTimer == nil {
		t.Error("Expected first timer to be set")
	}

	// Second copy should cancel first timer
	err = manager.Copy("second content", 500*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	if manager.clearTimer == firstTimer {
		t.Error("Expected new timer, got same timer")
	}
}

func TestClear(t *testing.T) {
	manager := New()

	// Set up a timer first
	err := manager.Copy("content to clear", 500*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	if manager.clearTimer == nil {
		t.Error("Expected timer to be set before clear")
	}

	// Clear should stop the timer
	err = manager.Clear()
	if err != nil {
		t.Errorf("Clear() failed: %v", err)
	}

	if manager.clearTimer != nil {
		t.Error("Expected timer to be nil after clear")
	}
}

func TestStopAutoClearing(t *testing.T) {
	manager := New()

	// Set up a timer
	err := manager.Copy("content", 500*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	if manager.clearTimer == nil {
		t.Error("Expected timer to be set before stopping")
	}

	manager.StopAutoClearing()

	if manager.clearTimer != nil {
		t.Error("Expected timer to be nil after StopAutoClearing")
	}
}

func TestStopAutoClearingWithoutTimer(t *testing.T) {
	manager := New()

	// Should not panic when no timer is set
	manager.StopAutoClearing()

	if manager.clearTimer != nil {
		t.Error("Expected timer to remain nil")
	}
}

func TestCopyEmptyString(t *testing.T) {
	manager := New()

	err := manager.Copy("", 100*time.Millisecond)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	// Should still set up timer even for empty string
	if manager.clearTimer == nil {
		t.Error("Expected timer even for empty string")
	}
}

func TestMultipleCopiesRapidly(t *testing.T) {
	manager := New()

	// Copy multiple times rapidly
	for i := 0; i < 5; i++ {
		err := manager.Copy("content "+string(rune('A'+i)), 200*time.Millisecond)
		if err != nil {
			t.Skipf("Clipboard not available in test environment: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Should have only one timer at the end
	if manager.clearTimer == nil {
		t.Error("Expected timer after multiple rapid copies")
	}
}

func TestCopyAndGet(t *testing.T) {
	manager := New()

	testContent := "test clipboard content"

	// Copy content
	err := manager.Copy(testContent, 0)
	if err != nil {
		t.Skipf("Clipboard not available in test environment: %v", err)
	}

	// Get content back
	retrieved, err := manager.Get()
	if err != nil {
		t.Errorf("Get() failed: %v", err)
	}

	if retrieved != testContent {
		t.Errorf("Expected '%s', got '%s'", testContent, retrieved)
	}
}
