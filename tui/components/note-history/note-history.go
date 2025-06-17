package noteHistory

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
)

type Model struct {
	tb     table.Model
	width  int
	height int
}

func NewModel() Model {
	columns := []table.Column{
		{Title: "Create Time", Width: 20},
		{Title: "ID", Width: 10},
		{Title: "Content", Width: 10},
	}

	tb := table.New(
		table.WithColumns(columns),
	)
	tb.Focus()

	// ti.Width = 20
	return Model{
		tb: tb,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case defs.CurrentViewMsg:
		if msg != "note-history" {
			m.tb.Blur()
		} else {
			m.tb.Focus()
		}
	}
	var cmd tea.Cmd
	m.tb, cmd = m.tb.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	noteStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Width(m.width).
		Height(m.height)

	noteView := noteStyle.Render(m.tb.View())
	// // Fill vertical space above the note to push it to the bottom
	// above := m.height - lipgloss.Height(noteView)
	// if above < 0 {
	// 	above = 0
	// }
	return noteView
}
