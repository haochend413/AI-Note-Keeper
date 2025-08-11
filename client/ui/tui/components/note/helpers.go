package note

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// m.ti.Width = width - 4
	m.ti.SetHeight(height - 1)
	m.ti.SetWidth(width - 1)
}

// sendnote
// type NoteSendMsg = *defs.Note

func (m *Model) SendNoteCmd() tea.Cmd {
	content := m.ti.Value()
	if content == "" {
		return nil
	}
	m.ti.Reset()
	return func() tea.Msg {
		return &defs.Note{
			Content: content,
		}
	}
}

func (m *Model) SendTopicCmd() tea.Cmd {
	content := m.ti.Value()
	if content == "" {
		return nil
	}
	m.ti.Reset()
	return func() tea.Msg {
		return &defs.Topic{
			Topic: content,
		}
	}
}

func (m *Model) ToggleEditable() tea.Cmd {
	// m.ti.SetEditable(!m.ti.Editable)
	// println(m.ti.Editable)
	return nil
}

func (m *Model) UpdateDisplay(content string) {
	//ah, i see, so everything just reset to zero after i press the button...ok.
	//so it is not the oter part's problem, it's my frontend. nice.
	//it should only be called when I change the selected note
	m.ti.SetValue(content)
}
