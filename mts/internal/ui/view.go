package ui

import (
	"strings"

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
		tableBox = focusedStyle.Render(m.table.View())
	} else {
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

	var topicsDisplay string
	var topicInputBox string
	if m.focus == FocusTopics {
		topicInputBox = focusedStyle.Render(m.topicInput.View())
	} else {
		topicInputBox = baseStyle.Render(m.topicInput.View())
	}

	if m.app.HasCurrentNote() {
		if topics := m.app.CurrentNoteTopics(); len(topics) > 0 {
			var topicTags []string
			maxWidth := m.width/2 - 8
			currentWidth := 0
			for _, topic := range topics {
				tagText := topic.Topic
				tagWidth := len(tagText) + 4
				if currentWidth+tagWidth > maxWidth && len(topicTags) > 0 {
					topicTags = append(topicTags, "\n")
					currentWidth = 0
				}
				topicTags = append(topicTags, topicStyle.Render(tagText))
				currentWidth += tagWidth
			}
			topicsDisplay = strings.Join(topicTags, "")
		} else {
			topicsDisplay = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("No topics")
		}
	} else {
		topicsDisplay = "No note selected"
	}

	rightSide := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("✏️  Edit Note"),
		editBox,
		titleStyle.Render("🏷️  Topics"),
		baseStyle.Width(max(20, m.width/2-4)).Height(4).Render(topicsDisplay),
		titleStyle.Render("➕ Add Topics"),
		topicInputBox,
	)

	main := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	help := helpStyle.Render(
		"Tab: cycle focus • Enter: select/search/add-topic • /: search • Ctrl+N: new note (table only) • Ctrl+S: save • Ctrl+Q: sync DB • Del: delete • Ctrl+C: quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}
