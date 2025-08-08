package dailyui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/app/state"
)

func StartDailyUI(appState *state.AppState) {
	p := tea.NewProgram(NewModel(appState))
	//model, error
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
