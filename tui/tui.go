package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/haochend413/mantis/tui/components/note"
	"github.com/haochend413/mantis/tui/keybindings"
)

type Model struct {
	// keybindings *keybindings.GlobalKeyMap
	noteModel note.Model
	width     int
	height    int
}

func NewModel() Model {
	return Model{
		noteModel: note.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// note update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.GlobalKeys.QuitApp):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.noteModel.SetSize(msg.Width, msg.Height/5)
	}
	m.noteModel, cmd = m.noteModel.Update(msg)
	return m, cmd
}

// Overall View management: positioning the views
func (m Model) View() string {
	noteView := m.noteModel.View()

	// Place the note at the bottom of the parent area
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Bottom,
		noteView,
	)
}
