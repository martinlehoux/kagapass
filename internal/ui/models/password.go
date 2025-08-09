package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/types"
)

// PasswordModel handles the password input screen
type PasswordModel struct {
	database     *types.Database
	password     string
	errorMessage string
	width        int
	height       int
	attempts     int
	maxAttempts  int
}

// NewPasswordModel creates a new password input model
func NewPasswordModel() *PasswordModel {
	return &PasswordModel{
		password:    "",
		attempts:    0,
		maxAttempts: 3,
	}
}

// SetDatabase sets the database to unlock
func (m *PasswordModel) SetDatabase(db *types.Database) {
	m.database = db
	m.password = ""
	m.errorMessage = ""
	m.attempts = 0
}

// Update implements tea.Model
func (m *PasswordModel) Update(msg tea.Msg) (*PasswordModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.password != "" {
				// Try to unlock the database
				return m, func() tea.Msg {
					return UnlockDatabaseMsg{
						Database: m.database,
						Password: m.password,
					}
				}
			}
		case "esc":
			// Return to file selection
			return m, func() tea.Msg {
				return SwitchScreenMsg{
					Screen: FileSelectionScreen,
				}
			}
		case "backspace":
			if len(m.password) > 0 {
				m.password = m.password[:len(m.password)-1]
			}
		case "ctrl+l":
			m.password = ""
		default:
			// Handle regular typing (but hide the characters)
			if len(msg.String()) == 1 {
				m.password += msg.String()
			}
		}
	case DatabaseUnlockResultMsg:
		if msg.Success {
			// Success! Switch to search screen
			return m, func() tea.Msg {
				return SwitchScreenMsg{
					Screen:   MainSearchScreen,
					Database: m.database,
					Entries:  msg.Entries,
				}
			}
		} else {
			// Failed to unlock
			m.attempts++
			m.password = ""
			m.errorMessage = msg.Error

			if m.attempts >= m.maxAttempts {
				// Too many attempts, return to file selection
				return m, func() tea.Msg {
					return SwitchScreenMsg{
						Screen: FileSelectionScreen,
					}
				}
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m *PasswordModel) View() string {
	var b strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Enter Master Password")

	b.WriteString(title + "\n\n")

	// Database info
	if m.database != nil {
		b.WriteString("Database: " + m.database.Name + "\n")
		if m.database.Path != "" {
			pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
			b.WriteString(pathStyle.Render("Path: "+m.database.Path) + "\n")
		}
		b.WriteString("\n")
	}

	// Error message
	if m.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		b.WriteString(errorStyle.Render("Error: "+m.errorMessage) + "\n\n")
	}

	// Show attempts remaining
	if m.attempts > 0 {
		attemptsLeft := m.maxAttempts - m.attempts
		warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD23F"))
		b.WriteString(warningStyle.Render("Incorrect password. Attempts remaining: ") +
			warningStyle.Render(string(rune('0'+attemptsLeft))) + "\n\n")
	}

	// Password input
	b.WriteString("Master Password:\n")

	// Show masked password
	maskedPassword := strings.Repeat("•", len(m.password))
	if len(m.password) == 0 {
		maskedPassword = "(empty)"
	}

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#2A2A2A")).
		Padding(0, 1).
		Width(30)

	b.WriteString(inputStyle.Render(maskedPassword) + "\n\n")

	// Footer
	footer := "[Enter] Unlock  [Esc] Back  [Ctrl+L] Clear"
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(footer))

	return b.String()
}

// UnlockDatabaseMsg is sent to attempt database unlocking
type UnlockDatabaseMsg struct {
	Database *types.Database
	Password string
}

// DatabaseUnlockResultMsg is sent with the result of database unlocking
type DatabaseUnlockResultMsg struct {
	Success bool
	Error   string
	Entries []types.Entry
}
