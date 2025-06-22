package tasklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) ToggleSuccess() tea.Cmd {
	return func() tea.Msg {
		m.TaskList[m.Index].Success = !m.TaskList[m.Index].Success
		return defs.TaskSucMsg{}
	}
}
