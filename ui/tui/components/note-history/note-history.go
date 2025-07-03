package noteHistory

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
	tui_defs "github.com/haochend413/mantis/defs/tui-defs"
)

type Model struct {
	tb      table.Model
	width   int
	height  int
	focus   bool
	context tui_defs.Context
}

func NewModel() Model {
	columns := []table.Column{
		{Title: "Create Time", Width: 20},
		{Title: "ID", Width: 5},
		{Title: "Content", Width: 10},
		{Title: "Topics", Width: 25},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{
		tb:      t,
		context: tui_defs.Default,
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

	return m.tb.View()
}
