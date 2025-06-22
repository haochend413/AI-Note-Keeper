package dailyui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func StartTui() {
	p := tea.NewProgram(NewModel())
	//model, error
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
