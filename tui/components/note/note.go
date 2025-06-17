package note

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
)

type Model struct {
	ti     textarea.Model
	width  int
	height int
}

func NewModel() Model {
	ti := textarea.New()
	ti.Placeholder = "Note"
	// ti.Focus()
	ti.CharLimit = 200
	// ti.Width = 20
	return Model{
		ti: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	// this is buggy; probably should not do that: use the mother component to handle everything, and even triggering the lower-level updates;
	switch msg := msg.(type) {
	case defs.CurrentViewMsg:
		if msg != "note" {
			m.ti.Blur()
		} else {
			m.ti.Focus()
		}
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd

}

func (m Model) View() string {
	noteStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Width(m.width).
		Height(m.height)

	noteView := noteStyle.Render(m.ti.View())
	// // Fill vertical space above the note to push it to the bottom
	// above := m.height - lipgloss.Height(noteView)
	// if above < 0 {
	// 	above = 0
	// }
	return noteView
}
