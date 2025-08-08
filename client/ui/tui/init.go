package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/app/state"
)

func StartTui(appState *state.AppState) {
	p := tea.NewProgram(NewModel(appState))
	//model, error
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
