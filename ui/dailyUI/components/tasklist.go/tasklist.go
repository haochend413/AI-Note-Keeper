package tasklist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/defs"
	"github.com/haochend413/mantis/ui/dailyUI/keybindings"
)

type Model struct {
	TaskList  []*defs.DailyTask
	Index     int
	Data      []*defs.DailyTask
	DBManager *dbcontroller.DBManager
}

func taskListView(m Model) string {
	var out string
	for _, task := range m.TaskList {
		out += checkbox(task.Task, task.Success) + "\n"
	}
	return mainStyle.Render(out)
}

// init to be emoty
func NewModel() Model {
	return Model{
		TaskList: []*defs.DailyTask{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			//pass data back to db
			m.DBManager.RefreshDaily(m.Data)
			return m, tea.Quit
		case key.Matches(msg, keybindings.DailyKeys.ToggleSuccess):
			//pass data back to db
			return m, m.ToggleSuccess()
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	return taskListView(m)
}
