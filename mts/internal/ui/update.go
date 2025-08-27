package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/bubbles/table"
)

// Update handles UI events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.statusBar.SetWidth(m.width)

		tableWidth := m.width/2 - 4
		editWidth := m.width/2 - 4
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Time", Width: 16},
			{Title: "Content", Width: max(10, tableWidth-70)},
			{Title: "Topics", Width: 20},
		}
		m.table.SetColumns(columns)
		if m.focus == FocusSearch {
			m.table.SetHeight(m.height - 10)
		} else {
			m.table.SetHeight(m.height - 8)
		}
		m.textarea.SetWidth(max(20, editWidth))
		m.textarea.SetHeight(max(5, m.height/2-6))
		m.searchInput.Width = tableWidth - 25
		m.topicInput.Width = max(20, editWidth)
		topicColumns := []table.Column{
			{Title: "Topic", Width: max(20, editWidth-4)},
		}
		m.topicsTable.SetColumns(topicColumns)
		m.topicsTable.SetWidth(max(20, editWidth))
		m.topicsTable.SetHeight(4)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.app.SyncWithDatabase()
			return m, tea.Quit
		case "tab":
			switch m.focus {
			case FocusSearch:
				m.focus = FocusTable
				m.searchInput.SetValue("")
				m.table.Focus()
				m.searchInput.Blur()
				m.topicInput.Blur()
				m.topicsTable.Blur()
				// m.table.SetHeight(m.height - 8)
			case FocusTable:
				m.focus = FocusEdit
				m.table.Blur()
				m.textarea.Focus()
				m.topicInput.Blur()
				m.topicsTable.Blur()
			case FocusEdit:
				m.app.SaveCurrentNote(m.textarea.Value())
				m.updateTable(Default)
				m.focus = FocusTopics
				m.textarea.Blur()
				m.topicInput.Focus()
				m.topicsTable.Focus()
			case FocusTopics:
				m.focus = FocusTable
				m.topicInput.Blur()
				m.topicsTable.Blur()
				m.table.Focus()
			}
			m.updateStatus()

		case "enter":
			switch m.focus {
			case FocusSearch:
				m.app.SearchNotes(m.searchInput.Value())
				m.focus = FocusTable
				m.table.Focus()
				m.searchInput.Blur()
				m.updateTable(Default)
			case FocusTable:
				m.app.SelectCurrentNote(m.table.Cursor())
				m.textarea.SetValue(m.app.CurrentNoteContent())
				m.updateTopicsTable()
			case FocusTopics:
				m.app.AddTopicsToCurrentNote(m.topicInput.Value())
				m.topicInput.SetValue("")
				m.updateTopicsTable()
			}
			m.updateStatus()

		case "ctrl+s":
			if m.focus == FocusEdit {
				m.app.SaveCurrentNote(m.textarea.Value())
				m.updateTable(Default)
				m.focus = FocusTable
				m.table.Focus()
				m.textarea.Blur()
				m.searchInput.Blur()
				m.topicInput.Blur()
				m.topicsTable.Blur()
				m.updateStatus()

			}
		case "ctrl+n", "n":
			if m.focus == FocusTable {
				m.app.CreateNewNote()
				m.table.Focus()
				m.textarea.SetValue(m.app.CurrentNoteContent())
				m.updateTable(Default)
				m.table.SetCursor(len(m.app.FilteredNotesList) - 1) // Set cursor to the last note
				m.updateStatus()
				//set focus to edit
				m.focus = FocusEdit
				m.table.Blur()
				m.textarea.Focus()
				m.topicInput.Blur()
				m.topicsTable.Blur()
				return m, nil
			}

		case "ctrl+q":
			m.app.SyncWithDatabase()
			// print("1111111")
			m.updateTable(Default)
			m.updateTopicsTable()
			m.updateFullTopicTable()
			m.app.UpdateRecentNotes()
			m.updateStatus()
			return m, nil
		case "backspace":
			switch m.focus {
			case FocusTable:
				m.app.DeleteCurrentNote(uint(m.table.Cursor()))
				// Adjust cursor to a valid position
				m.updateTable(Default)
				if len(m.app.FilteredNotesList) > 0 {
					newCursor := m.table.Cursor()
					if newCursor >= len(m.app.FilteredNotesList) {
						newCursor = len(m.app.FilteredNotesList) - 1
					}
					m.table.SetCursor(newCursor)
					m.app.SelectCurrentNote(newCursor)
				} else {
					m.app.SelectCurrentNote(0)
				}

				m.textarea.SetValue(m.app.CurrentNoteContent())
				m.updateTopicsTable()
				m.updateFullTopicTable()
			case FocusTopics:
				if m.app.HasCurrentNote() && len(m.app.CurrentNoteTopics()) > 0 {
					// Remove the selected topic from the current note
					topics := m.app.CurrentNoteTopics()
					cursor := m.topicsTable.Cursor()
					if cursor < len(topics) {
						m.app.RemoveTopicFromCurrentNote(topics[cursor].Topic)
						m.updateTopicsTable()
						// Adjust cursor if necessary
						if len(m.app.CurrentNoteTopics()) > 0 && cursor >= len(m.app.CurrentNoteTopics()) {
							m.topicsTable.SetCursor(len(m.app.CurrentNoteTopics()) - 1)
						}
					}
				}
				m.updateTable(Default) // Update note table to reflect topic changes
			}
		case "s":
			if m.focus == FocusTable {
				m.focus = FocusSearch
				m.searchInput.Focus()
				m.table.Blur()
				// m.table.SetHeight(m.height - 10)
				// m.searchInput.SetValue("")
				return m, nil // end the loop;
			}
		case "R":
			if m.focus == FocusTable {
				m.NoteSelector = Recent
				m.updateTable(Recent)
				m.updateStatus()
			}
			// return m, nil
		case "A":
			if m.focus == FocusTable {
				m.NoteSelector = Default
				m.updateTable(Default)
				m.updateStatus()

			}
			// return m, nil
			// case "e":
			// 	if m.focus == FocusTable {
			// 		m.app.HighlightCurrentNote()
			// 	}
			// 	return m, nil
		}

	case table.MoveSelectMsg:
		// tea.Printf("%s", "hihihihihihihi")
		switch m.focus {
		case FocusTable:
			m.app.SelectCurrentNote(m.table.Cursor())
			m.textarea.SetValue(m.app.CurrentNoteContent())
			m.updateTopicsTable()
		case FocusTopics:
			// No action needed for topic table selection, but ensure cursor is valid
			if len(m.app.CurrentNoteTopics()) > 0 && m.topicsTable.Cursor() >= len(m.app.CurrentNoteTopics()) {
				m.topicsTable.SetCursor(len(m.app.CurrentNoteTopics()) - 1)
			}
		}
		m.updateStatus()
	}

	switch m.focus {
	case FocusSearch:
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	case FocusTable:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	case FocusEdit:
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	case FocusTopics:
		var cmdTopicInput, cmdTopicsTable tea.Cmd
		m.topicInput, cmdTopicInput = m.topicInput.Update(msg)
		m.topicsTable, cmdTopicsTable = m.topicsTable.Update(msg)
		cmds = append(cmds, cmdTopicInput, cmdTopicsTable)
	}

	return m, tea.Batch(cmds...)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
