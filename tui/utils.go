package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) switchFocusCmd() tea.Cmd {
	return func() tea.Msg {
		var s defs.CurrentViewMsg
		switch m.AppStatus.CurrentView {
		case "note":
			m.AppStatus.CurrentView = "note-history"
			s = "note-history"
			// return "note-history"
		case "note-history":
			// return "note"
			m.AppStatus.CurrentView = "note"
			s = "note"
		default:
			m.AppStatus.CurrentView = "note"
			s = "note"
		}
		return s
	}
}
