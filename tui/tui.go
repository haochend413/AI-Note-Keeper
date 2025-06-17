package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/defs"
	"github.com/haochend413/mantis/tui/components/note"
	noteHistory "github.com/haochend413/mantis/tui/components/note-history"
	"github.com/haochend413/mantis/tui/keybindings"
)

type ViewType string

type Model struct {
	// keybindings *keybindings.GlobalKeyMap
	noteModel    note.Model
	historyModel noteHistory.Model
	//db
	DB_Data   *defs.DB_Data
	DBManager *dbcontroller.DBManager
	//size
	width  int
	height int
	//track
	AppStatus *defs.AppStatus
}

func NewModel() Model {
	return Model{
		noteModel:    note.NewModel(),
		historyModel: noteHistory.NewModel(),
		DB_Data:      &defs.DB_Data{},
		DBManager:    &dbcontroller.DBManager{},
		AppStatus: &defs.AppStatus{
			CurrentView: "note",
		},
	}
}

type initMsg struct {
	dumb bool
}

func (m *Model) initScreen() tea.Msg {
	//init db
	// m.DBManager.InitManager()
	// m.DB_Data = m.DBManager.FetchAll()
	// m.historyModel.UpdateDisplay(*m.DB_Data)
	return initMsg{dumb: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.initScreen, tea.EnterAltScreen)
}

// note update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		noteCmd    tea.Cmd
		historyCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case initMsg:
		//on init, load db data
		m.DBManager.InitManager()
		m.DB_Data = m.DBManager.FetchAll()
		m.historyModel.UpdateDisplay(*m.DB_Data)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			//pass data back to db
			m.DBManager.RefreshAll(m.DB_Data)
			return m, tea.Quit
		case key.Matches(msg, keybindings.GlobalKeys.SwitchFocus):
			return m, m.switchFocusCmd()
		case m.AppStatus.CurrentView == "note":
			switch {
			case key.Matches(msg, keybindings.Notekeys.SendNote):
				//send note to db
				return m, m.noteModel.SendNoteCmd()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.noteModel.SetSize(msg.Width-4, msg.Height/5)
		m.historyModel.SetSize(msg.Width/3, msg.Height/3*2)
		return m, nil
	case defs.NoteSendMsg:
		//update history section;
		m.DB_Data.NoteData = append(m.DB_Data.NoteData, msg)
		//update table display
		m.historyModel.UpdateDisplay(*m.DB_Data)
		return m, nil
	}
	m.noteModel, noteCmd = m.noteModel.Update(msg)
	m.historyModel, historyCmd = m.historyModel.Update(msg)

	return m, tea.Batch(noteCmd, historyCmd)
}

// Overall View management: positioning the views
func (m Model) View() string {
	noteView := m.noteModel.View()
	historyView := m.historyModel.View()

	// Place the note at the bottom of the parent area
	return lipgloss.JoinVertical(
		lipgloss.Top,
		historyView,
		noteView,
	)
}
