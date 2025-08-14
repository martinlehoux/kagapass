package models

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/martinlehoux/kagapass/internal/ui/status"
	"github.com/martinlehoux/kagapass/internal/ui/style"
)

// PasswordModel handles the password input screen
type PasswordModel struct {
	// Commands
	unlockDatabase *UnlockDatabase
	exit           func()

	database types.Database
	password string
	status   status.Status
}

func NewPasswordModel(unlockDatabase *UnlockDatabase, exit func(), database types.Database) *PasswordModel {
	return &PasswordModel{
		unlockDatabase: unlockDatabase,
		exit:           exit,
		database:       database,
		password:       "",
		status:         status.Status{},
	}
}

// Update implements tea.Model
func (m *PasswordModel) Update(msg tea.Msg) (*PasswordModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.password != "" {
				return m, m.unlockDatabase.Handle(m.database, []byte(m.password))
			}
		case "esc":
			m.exit()
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
	case DatabaseUnlockFailed:
		log.Printf("Failed to unlock database: %s", msg.Error)
	}
	return m, nil
}

// View implements tea.Model
func (m *PasswordModel) View() string {
	var b strings.Builder

	b.WriteString(style.ViewTitle.Render("Enter Master Password") + "\n\n")

	// Database info
	b.WriteString("Database: " + m.database.Name + "\n")
	if m.database.Path != "" {
		pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
		b.WriteString(pathStyle.Render("Path: "+m.database.Path) + "\n")
	}
	b.WriteString("\n")
	b.WriteString(m.status.Render() + "\n\n")
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
