package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/clipboard"
	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/sahilm/fuzzy"
)

// SearchModel handles the main search interface
type SearchModel struct {
	searchInput     string
	entries         []types.Entry
	filteredItems   []fuzzy.Match
	cursor          int
	width           int
	height          int
	clipboardManager *clipboard.Manager
	statusMessage   string
	dbName          string
}

// NewSearchModel creates a new search model
func NewSearchModel() *SearchModel {
	return &SearchModel{
		searchInput:      "",
		entries:          []types.Entry{},
		cursor:           0,
		clipboardManager: clipboard.New(),
	}
}

// SetClipboardManager sets the clipboard manager
func (m *SearchModel) SetClipboardManager(clipManager *clipboard.Manager) {
	m.clipboardManager = clipManager
}

// SetDatabaseName sets the current database name for display
func (m *SearchModel) SetDatabaseName(name string) {
	m.dbName = name
}

// SetEntries updates the entries list and triggers search
func (m *SearchModel) SetEntries(entries []types.Entry) {
	m.entries = entries
	m.search()
}

// Update implements tea.Model
func (m *SearchModel) Update(msg tea.Msg) (*SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filteredItems)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				entryIndex := m.filteredItems[m.cursor].Index
				if entryIndex < len(m.entries) {
					entry := &m.entries[entryIndex]
					return m, func() tea.Msg {
						return SwitchScreenMsg{
							Screen: EntryDetailsScreen,
							Entry:  entry,
						}
					}
				}
			}
		case "ctrl+b":
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				entryIndex := m.filteredItems[m.cursor].Index
				if entryIndex < len(m.entries) {
					entry := m.entries[entryIndex]
					if m.clipboardManager != nil && entry.Username != "" {
						err := m.clipboardManager.Copy(entry.Username, 30*time.Second)
						if err != nil {
							m.statusMessage = "Failed to copy username"
						} else {
							m.statusMessage = "Username copied to clipboard (will clear in 30s)"
						}
					} else {
						m.statusMessage = "No username to copy"
					}
				}
			}
		case "ctrl+c":
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				entryIndex := m.filteredItems[m.cursor].Index
				if entryIndex < len(m.entries) {
					entry := m.entries[entryIndex]
					if m.clipboardManager != nil && entry.Password != "" {
						err := m.clipboardManager.Copy(entry.Password, 30*time.Second)
						if err != nil {
							m.statusMessage = "Failed to copy password"
						} else {
							m.statusMessage = "Password copied to clipboard (will clear in 30s)"
						}
					} else {
						m.statusMessage = "No password to copy"
					}
				}
			}
		case "ctrl+l":
			m.searchInput = ""
			m.search()
		case "backspace":
			if len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.search()
			}
		default:
			// Handle regular typing
			if len(msg.String()) == 1 {
				m.searchInput += msg.String()
				m.search()
			}
		}
	}
	return m, nil
}

// View implements tea.Model
func (m *SearchModel) View() string {
	var b strings.Builder

	// Header
	titleText := "KagaPass"
	if m.dbName != "" {
		titleText += " - " + m.dbName
	}
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(titleText)

	b.WriteString(title + "\n\n")

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

	// Search input
	searchLabel := "Search: "
	searchValue := m.searchInput + "_" // Add cursor
	b.WriteString(searchLabel + searchValue + "\n")
	b.WriteString(strings.Repeat("─", 60) + "\n\n")

	// Results
	maxResults := 10
	if len(m.filteredItems) == 0 {
		if m.searchInput == "" {
			b.WriteString("Start typing to search entries...\n")
		} else {
			b.WriteString("No entries found.\n")
		}
	} else {
		for i, match := range m.filteredItems {
			if i >= maxResults {
				break
			}

			cursor := " "
			if m.cursor == i {
				cursor = "▶"
			}

			if match.Index < len(m.entries) {
				entry := m.entries[match.Index]
				line := fmt.Sprintf("  %s %s", cursor, entry.Title)
				if entry.Group != "" {
					line += fmt.Sprintf(" (%s)", entry.Group)
				}

				if m.cursor == i {
					line = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render(line)
				}

				b.WriteString(line + "\n")
			}
		}
	}

	b.WriteString("\n")

	// Footer
	footer := "[Ctrl+B] Copy User  [Ctrl+C] Copy Pass  [Enter] Details  [Esc] Files"
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render(footer))

	return m.wrapInBox(b.String())
}

// search performs fuzzy search on entries
func (m *SearchModel) search() {
	if m.searchInput == "" {
		m.filteredItems = []fuzzy.Match{}
		m.cursor = 0
		return
	}

	// Create search targets
	targets := make([]string, len(m.entries))
	for i, entry := range m.entries {
		targets[i] = entry.Title
	}

	// Perform fuzzy search
	matches := fuzzy.Find(m.searchInput, targets)
	m.filteredItems = matches
	m.cursor = 0
}

// wrapInBox wraps content in a border box
func (m *SearchModel) wrapInBox(content string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(65)

	return boxStyle.Render(content)
}