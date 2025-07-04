package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) switchFocusCmd() tea.Cmd {

	var nextView defs.CurrentViewMsg
	switch m.AppStatus.CurrentView {
	case "note":
		nextView = "note-history"
	case "note-history":

		nextView = "note"
	default:
		nextView = "note"
	}
	//return a msg for subcomponents to update their views
	return func() tea.Msg {
		return nextView
	}

}
