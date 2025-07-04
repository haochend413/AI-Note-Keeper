package note

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

type Model struct {
	ti     textarea.Model
	width  int
	height int
	focus  bool
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()
	return t
}
func NewModel() Model {
	ti := newTextarea()
	return Model{
		ti: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// note update function
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	// this is buggy; probably should not do that: use the mother component to handle everything, and even triggering the lower-level updates;
	switch msg := msg.(type) {
	case defs.CurrentViewMsg:
		if msg == "note" {
			m.focus = true
			m.ti.Focus()
		} else {
			m.focus = false
			m.ti.Blur()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	// noteStyle := lipgloss.NewStyle().
	// 	Border(lipgloss.NormalBorder()).
	// 	Padding(1, 1).
	// 	Width(m.width).
	// 	Height(m.height)

	// if m.focus {
	// 	noteStyle = noteStyle.
	// 		BorderForeground(lipgloss.Color("48"))
	// } else {
	// 	noteStyle = noteStyle.
	// 		BorderForeground(lipgloss.Color("15"))
	// }

	// noteView := noteStyle.Render(m.ti.View())
	// // // Fill vertical space above the note to push it to the bottom
	// // above := m.height - lipgloss.Height(noteView)
	// // if above < 0 {
	// // 	above = 0
	// // }
	return m.ti.View()
}
