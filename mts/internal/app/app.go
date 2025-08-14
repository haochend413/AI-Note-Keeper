package app

import (
	"strings"

	"github.com/haochend413/mts/internal/db"
	"github.com/haochend413/mts/internal/models"
)

// App encapsulates application logic and state
type App struct {
	db            *db.DB
	notes         []models.Note
	FilteredNotes []models.Note
	currentNote   *models.Note
	pendingNotes  []*models.Note
	deletedNotes  []uint
}

// NewApp creates a new application instance
func NewApp(dbConn *db.DB) *App {
	app := &App{
		db:           dbConn,
		pendingNotes: make([]*models.Note, 0),
		deletedNotes: make([]uint, 0),
	}
	app.loadNotes()
	return app
}

// loadNotes loads notes from the database
func (a *App) loadNotes() {
	notes, err := a.db.SyncWithDatabase([]models.Note{}, nil, nil)
	if err != nil {
		// Log error but continue with empty notes
		a.notes = []models.Note{}
		a.FilteredNotes = []models.Note{}
		return
	}
	a.notes = notes
	a.FilteredNotes = notes
}

// SearchNotes filters notes based on a query
func (a *App) SearchNotes(query string) {
	if query == "" {
		a.FilteredNotes = a.notes
		return
	}
	query = strings.ToLower(query)
	var filtered []models.Note
	for _, note := range a.notes {
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
	a.FilteredNotes = filtered
}

// SelectCurrentNote sets the current note based on table cursor
func (a *App) SelectCurrentNote(cursor int) {
	if len(a.FilteredNotes) == 0 || cursor >= len(a.FilteredNotes) {
		a.currentNote = nil
		return
	}
	a.currentNote = &a.FilteredNotes[cursor]
}

// CurrentNoteContent returns the content of the current note
func (a *App) CurrentNoteContent() string {
	if a.currentNote == nil {
		return ""
	}
	return a.currentNote.Content
}

// CurrentNoteTopics returns the topics of the current note
func (a *App) CurrentNoteTopics() []*models.Topic {
	if a.currentNote == nil {
		return nil
	}
	return a.currentNote.Topics
}

// HasCurrentNote checks if a note is currently selected
func (a *App) HasCurrentNote() bool {
	return a.currentNote != nil
}

// SaveCurrentNote updates the current note's content in-memory
func (a *App) SaveCurrentNote(content string) {
	if a.currentNote == nil {
		return
	}
	a.currentNote.Content = content
	for i, note := range a.notes {
		if note.ID == a.currentNote.ID && note.ID != 0 {
			a.notes[i].Content = content
			break
		}
	}
}

// AddTopicsToCurrentNote adds topics to the current note
func (a *App) AddTopicsToCurrentNote(topicsText string) {
	if a.currentNote == nil {
		return
	}
	topicsText = strings.TrimSpace(topicsText)
	if topicsText == "" {
		return
	}
	topicNames := strings.Split(topicsText, ",")
	for _, topicName := range topicNames {
		topicName = strings.TrimSpace(topicName)
		if topicName == "" {
			continue
		}
		topic := &models.Topic{Topic: topicName}
		exists := false
		for _, existing := range a.currentNote.Topics {
			if existing.Topic == topic.Topic {
				exists = true
				break
			}
		}
		if !exists {
			a.currentNote.Topics = append(a.currentNote.Topics, topic)
			for i, note := range a.notes {
				if note.ID == a.currentNote.ID && note.ID != 0 {
					a.notes[i].Topics = a.currentNote.Topics
					break
				}
			}
		}
	}
}

// RemoveTopicFromCurrentNote removes a topic from the current note
func (a *App) RemoveTopicFromCurrentNote(topicToRemove string) {
	if a.currentNote == nil {
		return
	}
	var newTopics []*models.Topic
	for _, topic := range a.currentNote.Topics {
		if topic.Topic != topicToRemove {
			newTopics = append(newTopics, topic)
		}
	}
	a.currentNote.Topics = newTopics
	for i, note := range a.notes {
		if note.ID == a.currentNote.ID && note.ID != 0 {
			a.notes[i].Topics = newTopics
			break
		}
	}
}

// CreateNewNote creates a new pending note
func (a *App) CreateNewNote(textareaValue string) {
	content := strings.TrimSpace(textareaValue)
	if content == "" {
		content = "New note"
	}
	note := &models.Note{Content: content}
	a.pendingNotes = append(a.pendingNotes, note)
	a.notes = append(a.notes, *note)
	a.FilteredNotes = append(a.FilteredNotes, *note)
	a.currentNote = &a.FilteredNotes[len(a.FilteredNotes)-1]
}

// DeleteCurrentNote deletes the current note in-memory
func (a *App) DeleteCurrentNote() {
	if a.currentNote == nil || len(a.FilteredNotes) == 0 {
		return
	}
	// Track deletion for database sync if note was in DB
	if a.currentNote.ID != 0 && !a.isPendingNote(a.currentNote) {
		a.deletedNotes = append(a.deletedNotes, a.currentNote.ID)
	}
	if a.isPendingNote(a.currentNote) {
		var newPending []*models.Note
		for _, pn := range a.pendingNotes {
			if pn != a.currentNote {
				newPending = append(newPending, pn)
			}
		}
		a.pendingNotes = newPending
	}
	// Remove from notes
	var newNotes []models.Note
	for _, note := range a.notes {
		if note.ID != a.currentNote.ID || a.isPendingNote(&note) {
			newNotes = append(newNotes, note)
		}
	}
	a.notes = newNotes
	// Remove from FilteredNotes
	var newFiltered []models.Note
	for _, note := range a.FilteredNotes {
		if note.ID != a.currentNote.ID || a.isPendingNote(&note) {
			newFiltered = append(newFiltered, note)
		}
	}
	a.FilteredNotes = newFiltered
	a.currentNote = nil
}

// isPendingNote checks if a note is pending
func (a *App) isPendingNote(note *models.Note) bool {
	for _, pn := range a.pendingNotes {
		if pn == note {
			return true
		}
	}
	return false
}

// HasChanges checks if there are unsaved changes
func (a *App) HasChanges() bool {
	if len(a.pendingNotes) > 0 || len(a.deletedNotes) > 0 {
		return true
	}
	for _, note := range a.notes {
		if note.ID != 0 { // Only check notes in DB
			var dbNote models.Note
			if err := a.db.Conn.First(&dbNote, note.ID).Error; err == nil {
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

// SyncWithDatabase syncs in-memory changes to the database
func (a *App) SyncWithDatabase() {
	notes, err := a.db.SyncWithDatabase(a.notes, a.pendingNotes, a.deletedNotes)
	if err != nil {
		// Log error but continue
		return
	}
	a.notes = notes
	a.FilteredNotes = notes
	a.pendingNotes = []*models.Note{}
	a.deletedNotes = []uint{}
}
