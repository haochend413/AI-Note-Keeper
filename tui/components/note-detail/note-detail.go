package noteDetail

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/mantis/defs"
)

type Model struct {
	vp       viewport.Model
	renderer *glamour.TermRenderer
	width    int
	height   int
	focus    bool
}

func NewModel() Model {
	vp := viewport.New(1, 1) // placeholder
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)

	//init renderer & initial display
	const glamourGutter = 2
	glamourRenderWidth := 78 - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		log.Panic(err)
	}

	str, err := renderer.Render("")
	if err != nil {
		log.Panic(err)
	}
	vp.SetContent(str)
	return Model{
		vp:       vp,
		renderer: renderer,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// note update function
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case defs.CurrentViewMsg:
		if msg == "note-detail" {
			m.focus = true
		} else {
			m.focus = false
		}
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	detailStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Width(m.width).
		Height(m.height)

	if m.focus {
		detailStyle = detailStyle.
			BorderForeground(lipgloss.Color("48")).
			Foreground(lipgloss.Color("15"))

	} else {
		detailStyle = detailStyle.
			BorderForeground(lipgloss.Color("15")).
			Foreground(lipgloss.Color("7"))

	}

	detailView := detailStyle.Render(m.vp.View())
	// // Fill vertical space above the note to push it to the bottom
	// above := m.height - lipgloss.Height(noteView)
	// if above < 0 {
	// 	above = 0
	// }
	return detailView
}
