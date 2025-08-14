package clipboard

import (
	"context"
	"log"
	"time"

	"github.com/atotto/clipboard"
)

// Clipboard handles clipboard operations with auto-clearing.
type Clipboard struct {
	clearTimer *time.Timer
}

// New creates a new clipboard manager.
func New() *Clipboard {
	return &Clipboard{
		clearTimer: nil,
	}
}

// Copy copies text to clipboard and sets up auto-clearing.
func (m *Clipboard) Copy(text string, clearAfter time.Duration) error {
	// Copy to clipboard
	err := clipboard.WriteAll(text)
	if err != nil {
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
				err := clipboard.WriteAll("")
				if err != nil {
					log.Printf("Failed to clear clipboard: %v", err)
				}
			}
		})
	}

	return nil
}

// CopyWithContext copies text with context cancellation support.
func (m *Clipboard) CopyWithContext(ctx context.Context, text string, clearAfter time.Duration) error {
	err := m.Copy(text, clearAfter)
	if err != nil {
		return err
	}

	// Handle context cancellation
	if clearAfter > 0 {
		go func() {
			select {
			case <-ctx.Done():
				// Context cancelled, clear immediately
				if current, err := clipboard.ReadAll(); err == nil && current == text {
					err := clipboard.WriteAll("")
					if err != nil {
						log.Printf("Failed to clear clipboard: %v", err)
					}
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

// Clear immediately clears the clipboard.
func (m *Clipboard) Clear() error {
	if m.clearTimer != nil {
		m.clearTimer.Stop()
		m.clearTimer = nil
	}

	return clipboard.WriteAll("")
}

// Get retrieves current clipboard content.
func (m *Clipboard) Get() (string, error) {
	return clipboard.ReadAll()
}

// StopAutoClearing cancels any pending auto-clear timer.
func (m *Clipboard) StopAutoClearing() {
	if m.clearTimer != nil {
		m.clearTimer.Stop()
		m.clearTimer = nil
	}
}
