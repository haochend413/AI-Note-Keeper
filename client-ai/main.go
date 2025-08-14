package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Data structures
type DB_Data struct {
	NoteData      []*Note
	TopicData     []*Topic
	DailyTaskData []*DailyTask
}

type Note struct {
	gorm.Model
	Content string
	Topics  []*Topic `gorm:"many2many:note_topics;"`
}

type Topic struct {
	gorm.Model
	Topic string
	Notes []*Note `gorm:"many2many:note_topics;"`
}

type DailyTask struct {
	gorm.Model
	Task    string
	Success bool
}

// Focus states
type focusState int

const (
	focusTable focusState = iota
	focusEdit
	focusSearch
	focusTopics
)

// Model
type model struct {
	db            *gorm.DB
	table         table.Model
	textarea      textarea.Model
	searchInput   textinput.Model
	topicInput    textinput.Model
	notes         []Note
	filteredNotes []Note
	currentNote   *Note
	focus         focusState
	width         int
	height        int
	ready         bool
	pendingNotes  []*Note // Notes not yet saved
	deletedNotes  []uint  // IDs of notes marked for deletion
}

// Styles
var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("211")).
			Bold(true).
			Padding(0, 1)

	topicStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Margin(0, 1, 0, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0, 0, 2)
)

func initialModel() model {
	db, err := gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&Note{}, &Topic{}, &DailyTask{})

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

	m := model{
		db:           db,
		table:        t,
		textarea:     ta,
		searchInput:  ti,
		topicInput:   topicInput,
		focus:        focusSearch,
		pendingNotes: make([]*Note, 0),
		deletedNotes: make([]uint, 0),
	}
	m.loadNotes()
	return m
}

// Load notes from DB (called only at startup and after sync)
func (m *model) loadNotes() {
	var notes []Note
	m.db.Preload("Topics").Find(&notes)
	m.notes = notes
	m.filteredNotes = notes
	m.updateTable()
}

// Update table rows
func (m *model) updateTable() {
	rows := make([]table.Row, len(m.filteredNotes))
	for i, note := range m.filteredNotes {
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
		// Use a placeholder ID and timestamp for pending notes
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

// Search notes
func (m *model) searchNotes(query string) {
	if query == "" {
		m.filteredNotes = m.notes
		m.updateTable()
		return
	}
	query = strings.ToLower(query)
	var filtered []Note
	for _, note := range m.notes {
		if strings.Contains(strings.ToLower(note.Content), query) {
			filtered = append(filtered, note)
			continue
		}
		for _, topic := range note.Topics {
			if strings.Contains(strings.ToLower(topic.Topic), query) {
				filtered = append(filtered, note)
				break
			}
		}
	}
	m.filteredNotes = filtered
	m.updateTable()
	m.table.SetCursor(0)
}

// Select current note
func (m *model) selectCurrentNote() {
	if len(m.filteredNotes) == 0 {
		m.currentNote = nil
		m.textarea.SetValue("")
		return
	}
	cursor := m.table.Cursor()
	if cursor < len(m.filteredNotes) {
		m.currentNote = &m.filteredNotes[cursor]
		m.textarea.SetValue(m.currentNote.Content)
	}
}

// Save current note to in-memory data
func (m *model) saveCurrentNote() {
	if m.currentNote == nil {
		return
	}
	m.currentNote.Content = m.textarea.Value()
	// Update m.notes if the note exists there
	for i, note := range m.notes {
		if note.ID == m.currentNote.ID && note.ID != 0 {
			m.notes[i].Content = m.currentNote.Content
			break
		}
	}
	m.updateTable()
}

// Check if note is pending
func (m *model) isPendingNote(note *Note) bool {
	for _, pn := range m.pendingNotes {
		if pn == note {
			return true
		}
	}
	return false
}

// Sync pending notes, updates, and deletions to DB
func (m *model) syncWithDatabase() {
	// Delete notes marked for deletion
	for _, noteID := range m.deletedNotes {
		m.db.Delete(&Note{}, noteID)
	}
	m.deletedNotes = []uint{} // Clear deleted notes

	// Save pending notes
	for _, note := range m.pendingNotes {
		note.Content = strings.TrimSpace(note.Content)
		if note.Content == "" {
			note.Content = "New note"
		}
		result := m.db.Create(note)
		if result.Error != nil {
			log.Printf("Error creating note in DB: %v", result.Error)
			continue
		}
		if len(note.Topics) > 0 {
			m.db.Model(note).Association("Topics").Append(note.Topics)
		}
	}
	m.pendingNotes = []*Note{} // Clear pending notes

	// Save updated notes
	for _, note := range m.notes {
		if note.ID != 0 { // Only save notes that were previously in DB
			m.db.Save(&note)
		}
	}

	m.loadNotes() // Reload from DB to get updated IDs and timestamps
}

// Add topics to current note
func (m *model) addTopicsToCurrentNote() {
	if m.currentNote == nil {
		return
	}
	topicsText := strings.TrimSpace(m.topicInput.Value())
	if topicsText == "" {
		return
	}
	topicNames := strings.Split(topicsText, ",")
	for _, topicName := range topicNames {
		topicName = strings.TrimSpace(topicName)
		if topicName == "" {
			continue
		}

		var topic Topic
		// For in-memory, create a new topic without DB interaction
		topic = Topic{Topic: topicName}
		exists := false
		for _, existing := range m.currentNote.Topics {
			if existing.Topic == topic.Topic {
				exists = true
				break
			}
		}
		if !exists {
			m.currentNote.Topics = append(m.currentNote.Topics, &topic)
			// Update m.notes if the note exists there
			for i, note := range m.notes {
				if note.ID == m.currentNote.ID && note.ID != 0 {
					m.notes[i].Topics = m.currentNote.Topics
					break
				}
			}
		}
	}
	m.topicInput.SetValue("")
	m.updateTable()
}

// Remove topic from current note
func (m *model) removeTopicFromCurrentNote(topicToRemove string) {
	if m.currentNote == nil {
		return
	}
	var newTopics []*Topic
	for _, topic := range m.currentNote.Topics {
		if topic.Topic != topicToRemove {
			newTopics = append(newTopics, topic)
		}
	}
	m.currentNote.Topics = newTopics
	// Update m.notes if the note exists there
	for i, note := range m.notes {
		if note.ID == m.currentNote.ID && note.ID != 0 {
			m.notes[i].Topics = newTopics
			break
		}
	}
	m.updateTable()
}

// Create new pending note
func (m *model) createNewNote() {
	content := strings.TrimSpace(m.textarea.Value())
	if content == "" {
		content = "New note"
	}
	note := &Note{Content: content}
	m.pendingNotes = append(m.pendingNotes, note)
	m.notes = append(m.notes, *note) // Add to m.notes for consistency
	m.filteredNotes = append(m.filteredNotes, *note)
	m.table.SetCursor(len(m.filteredNotes) - 1)
	m.selectCurrentNote()
	m.updateTable() // Ensure table is updated to show new note
}

// Delete current note
func (m *model) deleteCurrentNote() {
	if m.currentNote == nil || len(m.filteredNotes) == 0 {
		return
	}
	// Track deletion for database sync if note was in DB
	if m.currentNote.ID != 0 && !m.isPendingNote(m.currentNote) {
		m.deletedNotes = append(m.deletedNotes, m.currentNote.ID)
	}
	if m.isPendingNote(m.currentNote) {
		var newPending []*Note
		for _, pn := range m.pendingNotes {
			if pn != m.currentNote {
				newPending = append(newPending, pn)
			}
		}
		m.pendingNotes = newPending
	}
	// Remove from m.notes
	var newNotes []Note
	for _, note := range m.notes {
		if note.ID != m.currentNote.ID || m.isPendingNote(&note) {
			newNotes = append(newNotes, note)
		}
	}
	m.notes = newNotes
	// Remove from m.filteredNotes
	var newFiltered []Note
	for _, note := range m.filteredNotes {
		if note.ID != m.currentNote.ID || m.isPendingNote(&note) {
			newFiltered = append(newFiltered, note)
		}
	}
	m.filteredNotes = newFiltered
	m.currentNote = nil
	m.textarea.SetValue("")
	m.updateTable()
}

// Bubble Tea Init
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Bubble Tea Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case focusSearch:
				m.focus = focusTable
				m.table.Focus()
				m.searchInput.Blur()
				m.topicInput.Blur()
			case focusTable:
				m.focus = focusEdit
				m.table.Blur()
				m.textarea.Focus()
				m.topicInput.Blur()
			case focusEdit:
				m.focus = focusTopics
				m.textarea.Blur()
				m.topicInput.Focus()
			case focusTopics:
				m.focus = focusSearch
				m.topicInput.Blur()
				m.searchInput.Focus()
			}
		case "enter":
			if m.focus == focusSearch {
				m.searchNotes(m.searchInput.Value())
				m.focus = focusTable
				m.table.Focus()
				m.searchInput.Blur()
			} else if m.focus == focusTable {
				m.selectCurrentNote()
			} else if m.focus == focusTopics {
				m.addTopicsToCurrentNote()
			}
		case "ctrl+s":
			if m.focus == focusEdit {
				m.saveCurrentNote()
			}
		case "ctrl+n", "n":
			if m.focus == focusTable { // Only allow new notes in table view
				m.createNewNote()
				m.table.Focus() // Keep focus on table
			}
		case "ctrl+a":
			if m.focus == focusEdit {
				if m.currentNote == nil {
					m.createNewNote()
				} else {
					m.saveCurrentNote()
				}
			}
		case "ctrl+q":
			if len(m.pendingNotes) > 0 || len(m.deletedNotes) > 0 || m.hasUnsavedChanges() {
				m.syncWithDatabase()
			}
		case "delete":
			if m.focus == focusTable {
				m.deleteCurrentNote()
			} else if m.focus == focusTopics {
				m.topicInput.SetValue("")
			}
		case "/":
			if m.focus == focusTable {
				m.focus = focusSearch
				m.searchInput.Focus()
				m.table.Blur()
			}
		}
	}

	switch m.focus {
	case focusSearch:
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	case focusTable:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "up" || keyMsg.String() == "down" || keyMsg.String() == "j" || keyMsg.String() == "k" {
				m.selectCurrentNote()
			}
		}
	case focusEdit:
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	case focusTopics:
		m.topicInput, cmd = m.topicInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Check for unsaved changes in notes
func (m *model) hasUnsavedChanges() bool {
	for _, note := range m.notes {
		if note.ID != 0 { // Only check notes in DB
			var dbNote Note
			if err := m.db.First(&dbNote, note.ID).Error; err == nil {
				if dbNote.Content != note.Content || len(dbNote.Topics) != len(note.Topics) {
					return true
				}
				for i, topic := range dbNote.Topics {
					if i >= len(note.Topics) || topic.Topic != note.Topics[i].Topic {
						return true
					}
				}
			}
		}
	}
	return false
}

// Bubble Tea View
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var searchBox string
	if m.focus == focusSearch {
		searchBox = focusedStyle.Render(m.searchInput.View())
	} else {
		searchBox = baseStyle.Render(m.searchInput.View())
	}

	var tableBox string
	if m.focus == focusTable {
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
	if m.focus == focusEdit {
		editBox = focusedStyle.Render(m.textarea.View())
	} else {
		editBox = baseStyle.Render(m.textarea.View())
	}

	var topicsDisplay string
	var topicInputBox string
	if m.focus == focusTopics {
		topicInputBox = focusedStyle.Render(m.topicInput.View())
	} else {
		topicInputBox = baseStyle.Render(m.topicInput.View())
	}

	if m.currentNote != nil {
		if len(m.currentNote.Topics) > 0 {
			var topicTags []string
			maxWidth := m.width/2 - 8
			currentWidth := 0
			for _, topic := range m.currentNote.Topics {
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

// Helper
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Main
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
