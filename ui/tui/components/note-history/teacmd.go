package noteHistory

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
	tui_defs "github.com/haochend413/mantis/defs/tui-defs"
)

func (m *Model) SwitchContextCmd(nextContext tui_defs.Context) tea.Cmd {
	m.context = nextContext
	return func() tea.Msg {
		return defs.SwitchContextMsg{}
	}
}
