package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/types"
)

// FileSelectModel handles the file selection screen
type FileSelectModel struct {
	databases     types.DatabaseList
	cursor        int
	width         int
	height        int
	inputMode     bool
	inputText     string
	statusMessage string
}

// NewFileSelectModel creates a new file selection model
func NewFileSelectModel(databases types.DatabaseList) *FileSelectModel {
	return &FileSelectModel{
		databases: databases,
		cursor:    0,
	}
}

// Update implements tea.Model
func (m *FileSelectModel) Update(msg tea.Msg) (*FileSelectModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if !m.inputMode && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.inputMode && m.cursor < len(m.databases.Databases)-1 {
				m.cursor++
			}
		case "a":
			if !m.inputMode {
				m.inputMode = true
				m.inputText = ""
				m.statusMessage = "Enter path to KeePass database (.kdbx file):"
			}
		case "d":
			if !m.inputMode && len(m.databases.Databases) > 0 && m.cursor < len(m.databases.Databases) {
				// Remove selected database
				db := m.databases.Databases[m.cursor]
				m.databases.Databases = append(m.databases.Databases[:m.cursor], m.databases.Databases[m.cursor+1:]...)
				if m.cursor >= len(m.databases.Databases) && m.cursor > 0 {
					m.cursor--
				}
				m.statusMessage = fmt.Sprintf("Removed database: %s", db.Name)
				
				return m, func() tea.Msg {
					return UpdateDatabaseListMsg{DatabaseList: m.databases}
				}
			}
		case "esc":
			if m.inputMode {
				m.inputMode = false
				m.inputText = ""
				m.statusMessage = ""
			} else {
				return m, tea.Quit
			}
		case "enter":
			if m.inputMode {
				// Process the input
				return m.addDatabase()
			} else if len(m.databases.Databases) > 0 && m.cursor < len(m.databases.Databases) {
				selected := &m.databases.Databases[m.cursor]
				// Try to unlock with keyring first
				return m, func() tea.Msg {
					return TryKeyringUnlockMsg{Database: selected}
				}
			}
		case "backspace":
			if m.inputMode && len(m.inputText) > 0 {
				m.inputText = m.inputText[:len(m.inputText)-1]
			}
		default:
			// Handle input mode typing
			if m.inputMode && len(msg.String()) == 1 {
				m.inputText += msg.String()
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m *FileSelectModel) View() string {
	var b strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("KagaPass - Select Database")

	b.WriteString(title + "\n\n")

	// Show status message if any
	if m.statusMessage != "" {
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		b.WriteString(statusStyle.Render(m.statusMessage) + "\n\n")
	}

	// Handle input mode
	if m.inputMode {
		b.WriteString("Enter path to KeePass database (.kdbx file):\n\n")
		inputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#2A2A2A")).
			Padding(0, 1)
		b.WriteString(inputStyle.Render(m.inputText+"_") + "\n\n")
		b.WriteString("[Enter] Add  [Esc] Cancel\n")
		return m.wrapInBox(b.String())
	}

	if len(m.databases.Databases) == 0 {
		b.WriteString("No KeePass databases configured.\n\n")
		b.WriteString("Press 'a' to add a database file, 'Esc' to quit.\n")
		return m.wrapInBox(b.String())
	}

	b.WriteString("Select KeePass Database:\n\n")

	// Database list
	for i, db := range m.databases.Databases {
		cursor := " "
		if m.cursor == i {
			cursor = "▶"
		}

		name := db.Name
		if name == "" {
			name = fmt.Sprintf("Database %d", i+1)
		}

		line := fmt.Sprintf("  %s %s", cursor, name)
		if len(db.Path) > 0 {
			maxPathLen := 40
			path := db.Path
			if len(path) > maxPathLen {
				path = "..." + path[len(path)-maxPathLen+3:]
			}
			line += fmt.Sprintf("  (%s)", path)
		}

		if m.cursor == i {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render(line)
		}

		b.WriteString(line + "\n")
	}

	b.WriteString("\n")

	// Footer
	footer := "[Enter] Open  [Esc] Quit  [a] Add new file  [d] Remove"
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(footer))

	return m.wrapInBox(b.String())
}

// wrapInBox wraps content in a border box
func (m *FileSelectModel) wrapInBox(content string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(60)

	return boxStyle.Render(content)
}

// addDatabase validates and adds a new database to the list
func (m *FileSelectModel) addDatabase() (*FileSelectModel, tea.Cmd) {
	path := strings.TrimSpace(m.inputText)
	
	// Validate path
	if path == "" {
		m.statusMessage = "Path cannot be empty"
		return m, nil
	}
	
	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		m.statusMessage = "File does not exist"
		return m, nil
	}
	
	// Check if already in list
	for _, db := range m.databases.Databases {
		if db.Path == path {
			m.statusMessage = "Database already in list"
			return m, nil
		}
	}
	
	// Add database
	newDB := types.Database{
		Name:         filepath.Base(path),
		Path:         path,
		LastAccessed: time.Now(),
	}
	
	m.databases.Databases = append(m.databases.Databases, newDB)
	m.cursor = len(m.databases.Databases) - 1
	m.inputMode = false
	m.inputText = ""
	m.statusMessage = fmt.Sprintf("Added database: %s", newDB.Name)
	
	return m, func() tea.Msg {
		return UpdateDatabaseListMsg{DatabaseList: m.databases}
	}
}

// UpdateDatabaseListMsg is sent when the database list is modified
type UpdateDatabaseListMsg struct {
	DatabaseList types.DatabaseList
}

// TryKeyringUnlockMsg is sent to attempt unlocking with stored keyring password
type TryKeyringUnlockMsg struct {
	Database *types.Database
}