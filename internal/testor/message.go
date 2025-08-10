package testor

import tea "github.com/charmbracelet/bubbletea"

func KeyMsgRune(key rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}, Alt: false, Paste: false}
}
