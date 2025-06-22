package dailyui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/haochend413/mantis/app/state"
	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/defs"
	"github.com/haochend413/mantis/ui/dailyUI/components/tasklist.go"
	"github.com/haochend413/mantis/ui/tui/keybindings"
)

type ViewType string

type Model struct {
	// keybindings *keybindings.GlobalKeyMap
	taskModel tasklist.Model
	//db, passed as app state
	DB_Data   *defs.DB_Data
	DBManager *dbcontroller.DBManager
	//size
	// width  int
	// height int
}

func NewModel(appState *state.AppState) Model {
	return Model{
		taskModel: tasklist.NewModel(),
		DB_Data:   &appState.DB_Data,
		DBManager: appState.DBManager,
	}
}

func (m *Model) initScreen() tea.Msg {
	//init db
	// m.DBManager.InitManager()
	// m.DB_Data = m.DBManager.FetchAll()
	// m.historyModel.UpdateDisplay(*m.DB_Data)
	return defs.InitMsg{}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.initScreen)
}

// note update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taskCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case defs.InitMsg:
		//on init, load db data
		m.DBManager.InitManager()
		m.DB_Data = m.DBManager.FetchAll()
		m.taskModel.TaskList = m.DB_Data.DailyTaskData
		// return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			//pass data back to db
			m.DBManager.RefreshDaily(m.DB_Data.DailyTaskData)
			return m, tea.Quit

		}

	}
	m.taskModel.UpdateDisplay(m.DB_Data.DailyTaskData)
	m.taskModel, taskCmd = m.taskModel.Update(msg)

	return m, taskCmd
}

// Overall View management: positioning the views
func (m Model) View() string {
	return m.taskModel.View()

}
