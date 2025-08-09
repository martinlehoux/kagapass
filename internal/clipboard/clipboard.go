package clipboard

import (
	"context"
	"time"

	"github.com/atotto/clipboard"
)

// Manager handles clipboard operations with auto-clearing
type Manager struct {
	clearTimer *time.Timer
}

// New creates a new clipboard manager
func New() *Manager {
	return &Manager{}
}

// Copy copies text to clipboard and sets up auto-clearing
func (m *Manager) Copy(text string, clearAfter time.Duration) error {
	// Copy to clipboard
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	// Cancel any existing timer
	if m.clearTimer != nil {
		m.clearTimer.Stop()
	}

	// Set up auto-clear timer
	if clearAfter > 0 {
		m.clearTimer = time.AfterFunc(clearAfter, func() {
			// Check if clipboard still contains our text before clearing
			if current, err := clipboard.ReadAll(); err == nil && current == text {
				clipboard.WriteAll("") // Clear clipboard
			}
		})
	}

	return nil
}

// CopyWithContext copies text with context cancellation support
func (m *Manager) CopyWithContext(ctx context.Context, text string, clearAfter time.Duration) error {
	if err := m.Copy(text, clearAfter); err != nil {
		return err
	}

	// Handle context cancellation
	if clearAfter > 0 {
		go func() {
			select {
			case <-ctx.Done():
				// Context cancelled, clear immediately
				if current, err := clipboard.ReadAll(); err == nil && current == text {
					clipboard.WriteAll("")
				}
				if m.clearTimer != nil {
					m.clearTimer.Stop()
				}
			case <-time.After(clearAfter):
				// Timer completed naturally, nothing to do
			}
		}()
	}

	return nil
}

// Clear immediately clears the clipboard
func (m *Manager) Clear() error {
	if m.clearTimer != nil {
		m.clearTimer.Stop()
		m.clearTimer = nil
	}
	return clipboard.WriteAll("")
}

// Get retrieves current clipboard content
func (m *Manager) Get() (string, error) {
	return clipboard.ReadAll()
}

// StopAutoClearing cancels any pending auto-clear timer
func (m *Manager) StopAutoClearing() {
	if m.clearTimer != nil {
		m.clearTimer.Stop()
		m.clearTimer = nil
	}
}
