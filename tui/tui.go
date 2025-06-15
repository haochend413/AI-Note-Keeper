package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/haochend413/mantis/tui/components/note"
)

type Model struct {
	noteModel note.Model
}

func NewModel() Model {
	return Model{
		noteModel: note.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {

	// }
	var cmd tea.Cmd
	m.noteModel, cmd = m.noteModel.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.noteModel.View()
}
