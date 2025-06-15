package note

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ti textinput.Model
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Note"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return Model{
		ti: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// switch msg := msg.(type) {

	// }
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

// just for example
var noteStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(1, 2).
	Width(24)

func (m Model) View() string {
	return noteStyle.Render(m.ti.View())
}
