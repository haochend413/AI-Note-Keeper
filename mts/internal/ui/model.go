package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mts/internal/app"
	"github.com/haochend413/mts/internal/models"
)

// FocusState represents the current UI focus
type FocusState int

const (
	FocusTable FocusState = iota
	FocusEdit
	FocusSearch
	FocusTopics
)

// Model represents the Bubble Tea model
type Model struct {
	app         *app.App
	table       table.Model
	topicsTable table.Model
	textarea    textarea.Model
	searchInput textinput.Model
	topicInput  textinput.Model
	focus       FocusState
	width       int
	height      int
	ready       bool
}

// NewModel initializes a new UI model
func NewModel(application *app.App) Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Content", Width: 40},
		{Title: "Topics", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	topicColumns := []table.Column{
		{Title: "Topic", Width: 20},
	}
	//topic table
	tt := table.New(
		table.WithColumns(topicColumns),
		table.WithFocused(true),
		table.WithHeight(4),
	)

	ta := textarea.New()
	ta.Placeholder = "Edit note content..."
	ta.SetWidth(50)
	ta.SetHeight(10)

	ti := textinput.New()
	ti.Placeholder = "Search notes... (type to search, press Enter)"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	topicInput := textinput.New()
	topicInput.Placeholder = "Add topic (comma-separated)..."
	topicInput.CharLimit = 200
	topicInput.Width = 50

	m := Model{
		app:         application,
		table:       t,
		topicsTable: tt,
		textarea:    ta,
		searchInput: ti,
		topicInput:  topicInput,
		focus:       FocusSearch,
	}
	m.updateTable()
	m.updateTopicsTable()
	return m
}

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// updateTable updates the table rows based on filtered notes
func (m *Model) updateTable() {
	notes := make([]models.Note, 0, len(m.app.FilteredNotes))
	for _, note := range m.app.FilteredNotes {
		notes = append(notes, note)
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].CreatedAt.Before(notes[j].CreatedAt)
	})
	rows := make([]table.Row, len(notes))
	for i, note := range notes {
		topics := make([]string, len(note.Topics))
		for j, topic := range note.Topics {
			topics[j] = topic.Topic
		}
		topicsStr := strings.Join(topics, ", ")
		if len(topicsStr) > 18 {
			topicsStr = topicsStr[:15] + "..."
		}
		content := note.Content
		if len(content) > 38 {
			content = content[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", note.ID)
		timeStr := note.CreatedAt.Format("2006-01-02 15:04")
		if note.ID == 0 { // Pending note
			idStr = "P" // Indicate pending
			timeStr = time.Now().Format("2006-01-02 15:04")
		}
		rows[i] = table.Row{
			idStr,
			timeStr,
			content,
			topicsStr,
		}
	}
	m.table.SetRows(rows)
}

// updateTopicsTable updates the topics table rows based on the current note's topics
func (m *Model) updateTopicsTable() {
	rows := []table.Row{}
	if m.app.HasCurrentNote() {
		if topics := m.app.CurrentNoteTopics(); len(topics) > 0 {
			rows = make([]table.Row, len(topics))
			for i, topic := range topics {
				topicText := topic.Topic
				if len(topicText) > 18 {
					topicText = topicText[:15] + "..."
				}
				rows[i] = table.Row{topicText}
			}
		}
	}
	m.topicsTable.SetRows(rows)
}
