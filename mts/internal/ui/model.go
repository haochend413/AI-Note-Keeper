package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/haochend413/bubbles/statusbar"
	"github.com/haochend413/bubbles/table"
	"github.com/haochend413/bubbles/textarea"
	"github.com/haochend413/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mts/internal/app"
	"github.com/haochend413/mts/internal/models"
)

// FocusState represents the current UI focus
type Selector int

const (
	Default Selector = iota
	Recent
)

type FocusState int

const (
	FocusTable FocusState = iota
	FocusEdit
	FocusSearch
	FocusTopics
	FocusFullTopic
)

// Model represents the Bubble Tea model
type Model struct {
	app            *app.App
	NoteSelector   Selector
	table          table.Model
	fullTopicTable table.Model
	topicsTable    table.Model
	textarea       textarea.Model
	searchInput    textinput.Model
	topicInput     textinput.Model
	statusBar      statusbar.Model
	focus          FocusState
	width          int
	height         int
	ready          bool
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
	//set sb
	//Left: Context ; NoteID ; Last Update ; Version (frequency)
	//Right: Action ; Synced ? ; Time
	// Example usage with method chaining
	sb := statusbar.New(
		statusbar.WithHeight(1),
		statusbar.WithWidth(100),
		statusbar.WithLeftLen(4),
		statusbar.WithRightLen(3),
	)

	// Configure all left elements in sequence
	sb.GetLeft(0).SetValue("Context: Default").SetColors("0", "39").SetWidth(25)
	sb.GetLeft(1).SetValue("NoteID: -").SetColors("0", "45").SetWidth(15)
	sb.GetLeft(2).SetValue("Updated: Never").SetColors("0", "37").SetWidth(20)
	sb.GetLeft(3).SetValue("Version: 1.0").SetColors("0", "33").SetWidth(15)

	// Configure all right elements in sequence
	sb.GetRight(0).SetValue("Ready").SetColors("0", "46").SetWidth(12)
	sb.GetRight(1).SetValue("Not Synced").SetColors("0", "208").SetWidth(15)
	sb.GetRight(2).SetValue(time.Now().Format("15:04:05")).SetColors("0", "226").SetWidth(10)

	// You can also chain model methods
	sb.SetWidth(100).SetHeight(1)

	ta := textarea.New()
	ta.Placeholder = "Edit note content..."
	ta.SetWidth(50)
	ta.SetHeight(10)

	ti := textinput.New()
	ti.Placeholder = "Search notes... (type to search, press Enter)"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	fullTopicColumns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Topic", Width: 10},
	}

	ftt := table.New(

		table.WithColumns(fullTopicColumns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	topicInput := textinput.New()
	topicInput.Placeholder = "Add topic (comma-separated)..."
	topicInput.CharLimit = 200
	topicInput.Width = 50

	m := Model{
		app:            application,
		table:          t,
		topicsTable:    tt,
		textarea:       ta,
		searchInput:    ti,
		fullTopicTable: ftt,
		statusBar:      sb,
		topicInput:     topicInput,
		focus:          FocusTable,
	}
	m.updateTable(Default)
	m.updateTopicsTable()
	// print(len(m.app.Topics))
	m.updateFullTopicTable()
	return m
}

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) updateFullTopicTable() {
	topics := make([]models.Topic, 0, len(m.app.Topics))
	// print(len(m.app.Topics))
	for _, t := range m.app.Topics {
		topics = append(topics, *t)
	}
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].ID < topics[j].ID // Sort by ID (ascending)
		// Or sort alphabetically:
		// return topics[i].Topic < topics[j].Topic
	})
	rows := make([]table.Row, len(topics))
	for i, t := range topics {
		topicsStr := t.Topic
		if len(topicsStr) > 18 {
			topicsStr = topicsStr[:15] + "..."
		}
		idStr := fmt.Sprintf("%d", t.ID)
		rows[i] = table.Row{
			idStr,
			topicsStr,
		}
	}
	m.fullTopicTable.SetRows(rows)
	// print("aaaaa")
}

// updateTable updates the table rows based on the selector; it also updates the selector of the app;
func (m *Model) updateTable(s Selector) {
	m.NoteSelector = s
	var selectedNotes []*models.Note
	switch s {
	case Recent:
		selectedNotes = m.app.RecentNotes
	default:
		selectedNotes = m.app.FilteredNotesList
	}
	notes := make([]models.Note, 0, len(selectedNotes))
	for _, note := range selectedNotes {
		notes = append(notes, *note)
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ID < notes[j].ID
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
		timeStr := note.CreatedAt.Format("06-01-02 15:04")
		if note.ID == 0 { // Pending note
			idStr = "P" // Indicate pending
			timeStr = time.Now().Format("06-01-02 15:04")
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

// func (m *Model) updateStatusBar(s Selector) {

// }
