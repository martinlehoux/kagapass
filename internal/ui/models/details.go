package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/clipboard"
	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/martinlehoux/kagapass/internal/ui/style"
)

// DetailsModel handles the entry details screen
type DetailsModel struct {
	entry         types.Entry
	scroll        int
	clipboard     *clipboard.Clipboard
	statusMessage string
	showPassword  bool
}

// NewDetailsModel creates a new details model
func NewDetailsModel(clipboard *clipboard.Clipboard, entry types.Entry) *DetailsModel {
	return &DetailsModel{
		entry:         entry,
		scroll:        0,
		clipboard:     clipboard,
		statusMessage: "",
	}
}

// Update implements tea.Model
func (m *DetailsModel) Update(msg tea.Msg) (*DetailsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			// TODO: Implement proper scrolling based on content height
			m.scroll++
		case "ctrl+b":
			if m.clipboard != nil && m.entry.Username != "" {
				err := m.clipboard.Copy(m.entry.Username, 30*time.Second)
				if err != nil {
					m.statusMessage = "Failed to copy username"
				} else {
					m.statusMessage = "Username copied to clipboard (will clear in 30s)"
				}
			} else {
				m.statusMessage = "No username to copy"
			}
		case "ctrl+c":
			if m.clipboard != nil && m.entry.Password != "" {
				err := m.clipboard.Copy(m.entry.Password, 30*time.Second)
				if err != nil {
					m.statusMessage = "Failed to copy password"
				} else {
					m.statusMessage = "Password copied to clipboard (will clear in 30s)"
				}
			} else {
				m.statusMessage = "No password to copy"
			}
		case "ctrl+p":
			m.showPassword = !m.showPassword
			if m.showPassword {
				m.statusMessage = "Password revealed"
			} else {
				m.statusMessage = "Password hidden"
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m *DetailsModel) View() string {
	var b strings.Builder

	b.WriteString(style.ViewTitle.Render("Entry Details") + "\n\n")

	// Show status message if any
	if m.statusMessage != "" {
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#32D74B"))
		b.WriteString(statusStyle.Render(m.statusMessage) + "\n\n")
		// Clear status message after showing
		go func() {
			time.Sleep(3 * time.Second)
			m.statusMessage = ""
		}()
	}

	// Entry details
	b.WriteString(fmt.Sprintf("Title:    %s\n", m.entry.Title))
	b.WriteString(fmt.Sprintf("Username: %s\n", m.entry.Username))

	// Show password based on visibility toggle
	var passwordDisplay string
	if m.showPassword {
		passwordDisplay = m.entry.Password
	} else {
		passwordDisplay = strings.Repeat("*", 12)
	}
	b.WriteString(fmt.Sprintf("Password: %s\n", passwordDisplay))

	if m.entry.URL != "" {
		b.WriteString(fmt.Sprintf("URL:      %s\n", m.entry.URL))
	}

	if m.entry.Group != "" {
		b.WriteString(fmt.Sprintf("Group:    %s\n", m.entry.Group))
	}

	b.WriteString("\n")

	// Notes section
	if m.entry.Notes != "" {
		b.WriteString("Notes:\n")
		// TODO: Handle scrolling for long notes
		noteLines := strings.Split(m.entry.Notes, "\n")
		for _, line := range noteLines {
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}

	// Timestamps
	if !m.entry.Modified.IsZero() {
		b.WriteString(fmt.Sprintf("Modified: %s\n", m.entry.Modified.Format("2006-01-02 15:04:05")))
	}
	if !m.entry.Created.IsZero() {
		b.WriteString(fmt.Sprintf("Created:  %s\n", m.entry.Created.Format("2006-01-02 15:04:05")))
	}

	b.WriteString("\n")

	// Footer
	footer := "[Ctrl+B] Copy User  [Ctrl+C] Copy Pass  [Ctrl+P] Toggle Pass  [Esc] Back"
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(footer))

	return b.String()
}
