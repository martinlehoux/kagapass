package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagamigo/kcore"
	"github.com/martinlehoux/kagapass/internal/ui/models"
)

func main() {
	app, err := models.NewAppModel()
	kcore.Expect(err, "error initializing app")

	p := tea.NewProgram(app, tea.WithAltScreen())

	_, err = p.Run()
	kcore.Expect(err, "error running app")
}
