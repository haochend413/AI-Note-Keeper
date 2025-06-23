package tasklist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mantis/defs"
	"github.com/haochend413/mantis/ui/dailyUI/keybindings"
)

type Model struct {
	TaskList []*defs.DailyTask
	Index    int
}

// init to be emoty
func NewModel(data []*defs.DailyTask) Model {
	return Model{
		TaskList: data,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Index++
			if m.Index > len(m.TaskList) {
				m.Index = len(m.TaskList) - 1
			}
		case "up":
			m.Index--
			if m.Index < 0 {
				m.Index = 0
			}

		}
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			//pass data back to db

			return *m, tea.Quit
		case key.Matches(msg, keybindings.DailyKeys.ToggleSuccess):
			//pass data back to db
			return *m, m.ToggleSuccess()
		case key.Matches(msg, keybindings.DailyKeys.DeleteTask):
			//pass data back to db
			return *m, m.DeleteTask()
		}
		return *m, nil
	}
	return *m, nil
}

func (m *Model) View() string {
	return m.UpdateDisplay(m.TaskList)
}
