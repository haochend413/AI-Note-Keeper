package tasklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) ToggleSuccess() tea.Cmd {
	return func() tea.Msg {
		(*m.TaskList)[m.Index].Success = !(*m.TaskList)[m.Index].Success
		return defs.TaskSucMsg{}
	}
}

func (m *Model) DeleteTask() tea.Cmd {
	return func() tea.Msg {
		if len(*m.TaskList) == 0 {
			return defs.DeleteTaskMsg{}
		}
		*m.TaskList = append((*m.TaskList)[:m.Index], (*m.TaskList)[m.Index+1:]...)
		if m.Index >= len(*m.TaskList) && m.Index > 0 {
			m.Index--
		}
		if len(*m.TaskList) == 0 {
			m.Index = 0
		}
		return defs.DeleteTaskMsg{}
	}
}

var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)

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
