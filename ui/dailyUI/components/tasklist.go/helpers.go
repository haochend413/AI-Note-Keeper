package tasklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) ToggleSuccess() tea.Cmd {
	return func() tea.Msg {
		m.TaskList[m.Index].Success = !m.TaskList[m.Index].Success
		return defs.TaskSucMsg{}
	}
}

var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Background(lipgloss.Color("8"))

func (m *Model) UpdateDisplay(data []*defs.DailyTask) string {
	var out string
	for i, task := range data {
		line := checkbox(task.Task, task.Success)
		if i == m.Index {
			line = highlightStyle.Render(line)
		}
		out += line + "\n"
	}
	return mainStyle.Render(out)
}
