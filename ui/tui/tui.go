package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/haochend413/mantis/app/state"
	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/defs"
	tui_defs "github.com/haochend413/mantis/defs/tui-defs"
	"github.com/haochend413/mantis/ui/tui/components/note"
	noteDetail "github.com/haochend413/mantis/ui/tui/components/note-detail"
	noteHistory "github.com/haochend413/mantis/ui/tui/components/note-history"
	"github.com/haochend413/mantis/ui/tui/keybindings"
)

type ViewType string

type Model struct {
	// keybindings *keybindings.GlobalKeyMap
	noteModel    note.Model
	historyModel noteHistory.Model
	detailModal  noteDetail.Model
	//db
	DB_Data   *defs.DB_Data
	DBManager *dbcontroller.DBManager
	//size
	width  int
	height int
	//track
	AppStatus *tui_defs.AppStatus
}

func NewModel(appState *state.AppState) *Model {
	return &Model{
		noteModel:    note.NewModel(),
		historyModel: noteHistory.NewModel(),
		detailModal:  noteDetail.NewModel(),
		DB_Data:      appState.DB_Data,
		DBManager:    appState.DBManager,
		AppStatus: &tui_defs.AppStatus{
			//ok, more then "CurrentView", it should be understood as "NextView" :
			// what is the view that the model will switch into after hitting tab.
			CurrentView: "note-history",
		},
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
	return tea.Batch(m.initScreen, tea.EnterAltScreen)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		noteCmd    tea.Cmd
		historyCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case defs.InitMsg:
		//on init, load db data
		m.DBManager.InitManager()
		m.DB_Data = m.DBManager.FetchAll()
		m.historyModel.UpdateDisplay(*m.DB_Data)
		// return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			//pass data back to db
			m.DBManager.RefreshAll(m.DB_Data)
			return m, tea.Quit
		case key.Matches(msg, keybindings.GlobalKeys.FetchFromDB):
			//pass data back to db
			m.DBManager.RefreshAll(m.DB_Data)
			m.DB_Data = m.DBManager.FetchAll()
			m.historyModel.UpdateDisplay(*m.DB_Data)
			return m, nil
		case key.Matches(msg, keybindings.GlobalKeys.SwitchFocus):
			// Toggle view manually
			switch m.AppStatus.CurrentView {
			case "note":
				m.AppStatus.CurrentView = "note-history"
			case "note-history":
				m.AppStatus.CurrentView = "note"
			}
			return m, m.switchFocusCmd()
		//note view key bindings
		case m.AppStatus.CurrentView == "note":
			switch {
			case key.Matches(msg, keybindings.Notekeys.ToggleEditable):
				//send note to db
				return m, m.noteModel.ToggleEditable()
			case key.Matches(msg, keybindings.Notekeys.SendNote):
				//send note to db
				return m, m.noteModel.SendNoteCmd()
			case key.Matches(msg, keybindings.Notekeys.SendTopic):
				//send topic to db
				return m, m.noteModel.SendTopicCmd()
			}
		case m.AppStatus.CurrentView == "note-history":
			switch {
			case key.Matches(msg, keybindings.Historykeys.DayContext):
				return m, m.historyModel.SwitchContextCmd(tui_defs.Day)
			case key.Matches(msg, keybindings.Historykeys.MonthContext):
				return m, m.historyModel.SwitchContextCmd(tui_defs.Month)
			case key.Matches(msg, keybindings.Historykeys.WeekContext):
				return m, m.historyModel.SwitchContextCmd(tui_defs.Week)
			case key.Matches(msg, keybindings.Historykeys.DefaultContext):
				return m, m.historyModel.SwitchContextCmd(tui_defs.Default)
			case key.Matches(msg, keybindings.Historykeys.DeleteNote):
				return m, m.deleteNoteCmd()
			case key.Matches(msg, keybindings.Historykeys.EditNote):

				m.EditNote()
				return m, m.switchFocusCmd()
			}
		}

	case table.MoveSelectMsg:
		//update

		row := *msg.Row
		m.AppStatus.CurrentID, _ = strconv.Atoi(row[1])
		currentRow, _ := m.DBManager.FetchNoteFromID(m.AppStatus.CurrentID)

		// Update content only if row changed

		content := ""
		if currentRow != nil {
			content = currentRow.Content

		}

		m.noteModel.UpdateDisplay(content)
		m.detailModal.UpdateDisplay(content)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.noteModel.SetSize(msg.Width*2/3-4, msg.Height/3)
		m.historyModel.SetSize(msg.Width/3, msg.Height-2)
		m.detailModal.SetSize(msg.Width/3*2-5, msg.Height/3)
		return m, nil
	case defs.NoteSendMsg:
		//update history section;
		//.. in this case is lazy-sync a good idea ?
		// rethink the whole idea;
		m.DBManager.UpdateNote(strconv.Itoa(m.AppStatus.CurrentID), msg.Content)
		m.DB_Data.NoteData = append(m.DB_Data.NoteData, msg)
		//update table display
		m.historyModel.UpdateDisplay(*m.DB_Data)
		return m, nil
	case defs.TopicSendMsg:
		//update history section;
		//just first rewrite to database, and then fetch;

		m.DB_Data.TopicData = append(m.DB_Data.TopicData, msg)
		//update table display
		m.historyModel.UpdateDisplay(*m.DB_Data)
		return m, nil
	case defs.UpdateHistoryDisplay:
		m.historyModel.UpdateDisplay(*m.DB_Data)
		// case defs.DeleteNoteMsg:
		// 	m.historyModel.UpdateDisplay(*m.DB_Data)
	}

	m.noteModel, noteCmd = m.noteModel.Update(msg)
	m.historyModel, historyCmd = m.historyModel.Update(msg)

	// // Get current row
	// currentRow := m.historyModel.GetCurrentRowData()

	// // Only update if row changed (by comparing with last selected row)
	// rowChanged := false
	// if len(currentRow) > 0 && len(m.AppStatus.LastRowSelected) > 0 && len(currentRow) == len(m.AppStatus.LastRowSelected) {
	// 	// Check if IDs are different (assuming ID is in index 1)

	// 	if currentRow[1] != m.AppStatus.LastRowSelected[1] {
	// 		rowChanged = true
	// 	}
	// } else if len(currentRow) != len(m.AppStatus.LastRowSelected) {
	// 	// Different length means different rows
	// 	rowChanged = true
	// }

	//fetch current note

	return m, tea.Batch(noteCmd, historyCmd)
}

func (m Model) View() string {
	noteView := m.noteModel.View()
	historyView := m.historyModel.View()
	// detailView := m.detailModal.View()

	// Place the note view to the right of the history view (both at the top)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		historyView,
		noteView,
	)
}
