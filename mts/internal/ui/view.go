package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var searchBox string
	if m.focus == FocusSearch {
		searchBox = focusedStyle.Render(m.searchInput.View())
	} else {
		searchBox = baseStyle.Render(m.searchInput.View())
	}

	var tableBox string
	if m.focus == FocusTable {
		m.table.SetStyles(focusedTableStyle)
		tableBox = focusedStyle.Render(m.table.View())
	} else {
		m.table.SetStyles(baseTableStyle)
		tableBox = baseStyle.Render(m.table.View())
	}

	leftSide := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("🔍 Search"),
		searchBox,
		titleStyle.Render("📝 Notes"),
		tableBox,
	)

	var editBox string
	if m.focus == FocusEdit {
		editBox = focusedStyle.Render(m.textarea.View())
	} else {
		editBox = baseStyle.Render(m.textarea.View())
	}

	var topicsTableBox string
	if m.focus == FocusTopics {
		// When focused, use a minimal highlight style
		topicsTableBox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")).
			Width(max(20, m.width/2-4)).
			MaxHeight(4).
			Render(m.topicsTable.View())
	} else {
		// When not focused, use an even more minimal style
		topicsTableBox = lipgloss.NewStyle().
			Width(max(20, m.width/2-4)).
			MaxHeight(4).
			Render(m.topicsTable.View())
	}

	var topicInputBox string
	if m.focus == FocusTopics {
		topicInputBox = focusedStyle.Render(m.topicInput.View())
	} else {
		topicInputBox = baseStyle.Render(m.topicInput.View())
	}

	rightSide := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("✏️  Edit Note"),
		editBox,
		titleStyle.Render("🏷️  Topics"),
		simpleTopicsStyle.Render(topicsTableBox),
		titleStyle.Render("➕ Add Topics"),
		topicInputBox,
	)

	main := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	help := helpStyle.Render(
		"Tab: cycle focus • Enter: select/search/add-topic • /: search • Ctrl+N: new note (table only) • Ctrl+S: save • Ctrl+Q: sync DB • Del: delete note/topic • Ctrl+C: quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}
