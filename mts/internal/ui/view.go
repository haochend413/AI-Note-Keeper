package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate space for main content (everything except help and status bar)
	// Reserve 2 lines for help and status bar
	mainContentHeight := m.height - 3

	// Calculate table height based on this available space
	tableHeight := mainContentHeight - 6 // Default height when search is not visible

	var leftSide string
	if m.focus == FocusSearch {
		// When search is visible, make the table shorter
		tableHeight = mainContentHeight - 9 // Reduce height to accommodate search

		// Set the table height dynamically
		m.table.SetHeight(tableHeight)

		searchBox := focusedStyle.Render(m.searchInput.View())
		tableBox := m.renderTableBox()

		// Join with bottom alignment to ensure content grows from top
		leftSide = lipgloss.JoinVertical(lipgloss.Bottom,
			searchBox,
			tableBox,
		)
	} else {
		// When search is not visible, make the table taller
		m.table.SetHeight(tableHeight)

		tableBox := m.renderTableBox()

		// No search box, just render the table with bottom alignment
		leftSide = tableBox
	}

	var fullTopicTableBox string
	if m.focus == FocusFullTopic {
		m.fullTopicTable.SetStyles(focusedTableStyle)
		fullTopicTableBox = focusedStyle.Render(m.fullTopicTable.View())
	} else {
		m.fullTopicTable.SetStyles(baseTableStyle)
		fullTopicTableBox = baseStyle.Render(m.fullTopicTable.View())
	}

	realLeft := lipgloss.JoinHorizontal(lipgloss.Left,
		fullTopicTableBox,
		leftSide,
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
		editBox,
		titleStyle.Render("Topics"),
		simpleTopicsStyle.Render(topicsTableBox),
		titleStyle.Render("Add Topics"),
		topicInputBox,
	)

	// Create main content with fixed height to ensure bottom elements are pushed down
	mainContent := lipgloss.NewStyle().
		Height(mainContentHeight).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, realLeft, rightSide))

	help := helpStyle.Render(
		"Tab: cycle focus • Enter: select/search/add-topic • /: search • Ctrl+N: new note (table only) • Ctrl+S: save • Ctrl+Q: sync DB • Del: delete note/topic • Ctrl+C: quit",
	)

	// Render status bar without extra styling that might add space
	statusBarBox := m.statusBar.View()

	// Join everything with the main content taking up all available space except for help and status
	return lipgloss.JoinVertical(lipgloss.Top,
		mainContent,
		help,
		statusBarBox,
	)
}

func (m Model) renderTableBox() string {
	if m.focus == FocusTable {
		m.table.SetStyles(focusedTableStyle)
		return focusedStyle.Render(m.table.View())
	} else {
		m.table.SetStyles(baseTableStyle)
		return baseStyle.Render(m.table.View())
	}
}
