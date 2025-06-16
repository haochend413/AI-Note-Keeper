package note

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/models"
)

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// m.ti.Width = width - 4
}

// sendnote
type NoteSendMsg = *models.Note

func (m *Model) SendNoteCmd() tea.Cmd {
	content := m.ti.Value()
	if content == "" {
		return nil
	}
	m.ti.Reset()
	return func() tea.Msg {
		return &models.Note{
			Content: content,
		}
	}
}
