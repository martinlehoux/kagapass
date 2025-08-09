package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/martinlehoux/kagapass/internal/ui/style"
)

// FileSelectModel handles the file selection screen
type FileSelectModel struct {
	// Commands
	unlockDatabase func(database types.Database, password string) tea.Cmd

	databases     types.DatabaseList
	cursor        int
	databaseInput textinput.Model
	statusMessage string
}

// NewFileSelectModel creates a new file selection model
func NewFileSelectModel(databases types.DatabaseList, unlockDatabase func(database types.Database, password string) tea.Cmd) *FileSelectModel {
	return &FileSelectModel{
		unlockDatabase: unlockDatabase,
		databases:      databases,
		databaseInput:  textinput.New(),
		cursor:         0,
		statusMessage:  "",
	}
}

func (m *FileSelectModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m *FileSelectModel) Update(msg tea.Msg) (*FileSelectModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input mode first - takes priority over navigation
		if m.databaseInput.Focused() {
			switch msg.String() {
			case "esc":
				m.databaseInput.Blur()
				m.databaseInput.Reset()
				m.statusMessage = ""
			case "enter":
				return m.addDatabase()
			default:
				m.databaseInput, cmd = m.databaseInput.Update(msg)
				return m, cmd
			}
		} else {
			// Handle navigation mode
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.databases.Databases)-1 {
					m.cursor++
				}
			case "a":
				m.databaseInput.Focus()
				m.statusMessage = ""
			case "d":
				m, cmd = m.removeDatabase()
				return m, cmd
			case "enter":
				if len(m.databases.Databases) > 0 && m.cursor < len(m.databases.Databases) {
					return m, m.unlockDatabase(m.databases.Databases[m.cursor], "")
				}
			case "esc":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m *FileSelectModel) View() string {
	var b strings.Builder

	title := style.ViewTitle.Render("KagaPass - Select Database")
	b.WriteString(title + "\n\n")

	if m.statusMessage != "" {
		b.WriteString(style.StatusMessage.Render(m.statusMessage) + "\n\n")
	}

	if m.databaseInput.Focused() {
		b.WriteString("Enter path to KeePass database (.kdbx file):\n\n")
		b.WriteString(m.databaseInput.View() + "\n\n")
		b.WriteString("[Enter] Add  [Esc] Cancel\n")
		return b.String()
	}

	if len(m.databases.Databases) == 0 {
		b.WriteString("No KeePass databases configured.\n\n")
		b.WriteString("Press 'a' to add a database file, 'Esc' to quit.\n")
		return b.String()
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

	return b.String()
}

// addDatabase validates and adds a new database to the list
func (m *FileSelectModel) addDatabase() (*FileSelectModel, tea.Cmd) {
	path := strings.TrimSpace(m.databaseInput.Value())

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
	m.databaseInput.Blur()
	m.databaseInput.Reset()
	m.statusMessage = fmt.Sprintf("Added database: %s", newDB.Name)

	return m, func() tea.Msg {
		return UpdateDatabaseListMsg{DatabaseList: m.databases}
	}
}

func (m *FileSelectModel) removeDatabase() (*FileSelectModel, tea.Cmd) {
	if m.cursor < 0 || m.cursor >= len(m.databases.Databases) {
		return m, nil
	}

	deleted := m.databases.Databases[m.cursor]
	m.databases.Databases = append(m.databases.Databases[:m.cursor], m.databases.Databases[m.cursor+1:]...)
	m.cursor = max(0, m.cursor-1)
	m.statusMessage = fmt.Sprintf("Removed database: %s", deleted.Name)

	return m, func() tea.Msg {
		return UpdateDatabaseListMsg{DatabaseList: m.databases}
	}
}

// UpdateDatabaseListMsg is sent when the database list is modified
type UpdateDatabaseListMsg struct {
	DatabaseList types.DatabaseList
}
