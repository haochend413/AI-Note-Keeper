package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
		tableWidth := m.width/2 - 4
		editWidth := m.width/2 - 4
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Time", Width: 16},
			{Title: "Content", Width: max(10, tableWidth-45)},
			{Title: "Topics", Width: 20},
		}
		m.table.SetColumns(columns)
		m.table.SetHeight(m.height - 8)
		m.textarea.SetWidth(max(20, editWidth))
		m.textarea.SetHeight(max(5, m.height/2-6))
		m.searchInput.Width = max(20, tableWidth)
		m.topicInput.Width = max(20, editWidth)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.focus {
			case FocusSearch:
				m.focus = FocusTable
				m.table.Focus()
				m.searchInput.Blur()
				m.topicInput.Blur()
			case FocusTable:
				m.focus = FocusEdit
				m.table.Blur()
				m.textarea.Focus()
				m.topicInput.Blur()
			case FocusEdit:
				m.focus = FocusTopics
				m.textarea.Blur()
				m.topicInput.Focus()
			case FocusTopics:
				m.focus = FocusSearch
				m.topicInput.Blur()
				m.searchInput.Focus()
			}
		case "enter":
			switch m.focus {
			case FocusSearch:
				m.app.SearchNotes(m.searchInput.Value())
				m.focus = FocusTable
				m.table.Focus()
				m.searchInput.Blur()
				m.updateTable()
			case FocusTable:
				m.app.SelectCurrentNote(m.table.Cursor())
				m.textarea.SetValue(m.app.CurrentNoteContent())
			case FocusTopics:
				m.app.AddTopicsToCurrentNote(m.topicInput.Value())
				m.topicInput.SetValue("")
				m.updateTable()
			}
		case "ctrl+s":
			if m.focus == FocusEdit {
				m.app.SaveCurrentNote(m.textarea.Value())
				m.updateTable()
			}
		case "ctrl+n", "n":
			if m.focus == FocusTable {
				m.app.CreateNewNote(m.textarea.Value())
				m.table.Focus()
				m.textarea.SetValue(m.app.CurrentNoteContent())
				m.updateTable()
			}
		case "ctrl+a":
			if m.focus == FocusEdit {
				if !m.app.HasCurrentNote() {
					m.app.CreateNewNote(m.textarea.Value())
					m.textarea.SetValue(m.app.CurrentNoteContent())
				} else {
					m.app.SaveCurrentNote(m.textarea.Value())
				}
				m.updateTable()
			}
		case "ctrl+q":
			if m.app.HasChanges() {
				m.app.SyncWithDatabase()
				m.updateTable()
			}
		case "delete":
			switch m.focus {
			case FocusTable:
				m.app.DeleteCurrentNote()
				// m.textarea.SetValue(m.app.CurrentNoteContent())
				m.app.SelectCurrentNote(m.table.Cursor())
				m.textarea.SetValue(m.app.CurrentNoteContent())
				m.updateTable()
			case FocusTopics:
				m.topicInput.SetValue("")
			}
		case "/":
			if m.focus == FocusTable {
				m.focus = FocusSearch
				m.searchInput.Focus()
				m.table.Blur()
			}
		}
	}

	switch m.focus {
	case FocusSearch:
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	case FocusTable:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "up" || keyMsg.String() == "down" || keyMsg.String() == "j" || keyMsg.String() == "k" {
				m.app.SelectCurrentNote(m.table.Cursor())
				m.textarea.SetValue(m.app.CurrentNoteContent())
			}
		}
	case FocusEdit:
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	case FocusTopics:
		m.topicInput, cmd = m.topicInput.Update(msg)
		cmds = append(cmds, cmd)
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
