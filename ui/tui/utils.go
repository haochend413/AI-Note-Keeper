package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) switchFocusCmd() tea.Cmd {
	nextView := defs.CurrentViewMsg(m.AppStatus.CurrentView)
	//return a msg for subcomponents to update their views
	return func() tea.Msg {
		return nextView
	}
}

func (m *Model) deleteNoteCmd() tea.Cmd {
	// This is O(n). Remove from db and make the sync process into daemon is O(1); we can do that later.Optimization is crucial;
	selectid64, _ := strconv.ParseUint(m.historyModel.GetCurrentRowData()[1], 10, 64)
	selectid := uint(selectid64)
	newNotes := m.DB_Data.NoteData[:0]
	for _, note := range m.DB_Data.NoteData {
		if note.ID != selectid {
			newNotes = append(newNotes, note)
		}
	}
	m.DB_Data.NoteData = newNotes
	return func() tea.Msg {
		return defs.DeleteNoteMsg{}
	}
}
