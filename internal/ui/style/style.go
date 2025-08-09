package style

import "github.com/charmbracelet/lipgloss"

var ViewTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(0, 1)

var StatusMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
