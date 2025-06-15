package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func StartTui() {
	p := tea.NewProgram(NewModel())
	if err, _ := p.Run(); err != nil {
		panic(err)
	}
}
