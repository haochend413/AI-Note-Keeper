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
	focus  bool
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
		if msg == "note-history" {
			m.focus = true
			m.tb.Focus()
		} else {
			m.focus = false
			m.tb.Blur()
		}
	}
	var cmd tea.Cmd
	m.tb, cmd = m.tb.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	historyStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Width(m.width).
		Height(m.height)

	if m.focus {
		historyStyle = historyStyle.
			BorderForeground(lipgloss.Color("48")).
			Foreground(lipgloss.Color("15"))

	} else {
		historyStyle = historyStyle.
			BorderForeground(lipgloss.Color("15")).
			Foreground(lipgloss.Color("7"))

	}

	historyView := historyStyle.Render(m.tb.View())
	// // Fill vertical space above the note to push it to the bottom
	// above := m.height - lipgloss.Height(noteView)
	// if above < 0 {
	// 	above = 0
	// }
	return historyView
}
