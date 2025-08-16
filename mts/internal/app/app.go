package app

import (
	"sort"
	"strings"
	"sync"

	"github.com/haochend413/mts/internal/db"
	"github.com/haochend413/mts/internal/models"
)

// App encapsulates application logic and state
type App struct {
	db            *db.DB
	notes         map[uint]models.Note // Notes by ID (or temp ID for pending)
	FilteredNotes map[uint]models.Note // Filtered notes by ID
	currentNote   *models.Note
	pendingNotes  map[*models.Note]uint // Map pending notes to their temp IDs
	deletedNotes  map[uint]struct{}     // IDs of deleted notes
	nextTempID    uint                  // For generating temporary IDs for pending notes
	mutex         sync.Mutex            // Protect concurrent access to maps
}

// NewApp creates a new application instance
func NewApp(dbConn *db.DB) *App {
	app := &App{
		db:            dbConn,
		notes:         make(map[uint]models.Note),
		FilteredNotes: make(map[uint]models.Note),
		pendingNotes:  make(map[*models.Note]uint),
		deletedNotes:  make(map[uint]struct{}),
		nextTempID:    1, // Start temporary IDs from 1
	}
	app.loadNotes()
	return app
}

// loadNotes loads notes from the database
func (a *App) loadNotes() {
	notes, err := a.db.SyncWithDatabase([]models.Note{}, nil, nil)
	if err != nil {
		// Log error but continue with empty notes
		a.notes = make(map[uint]models.Note)
		a.FilteredNotes = make(map[uint]models.Note)
		return
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.notes = make(map[uint]models.Note, len(notes))
	a.FilteredNotes = make(map[uint]models.Note, len(notes))
	maxID := uint(0)
	for _, note := range notes {
		a.notes[note.ID] = note
		a.FilteredNotes[note.ID] = note
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	a.nextTempID = maxID + 1
}

// SearchNotes filters notes based on a query
func (a *App) SearchNotes(query string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if query == "" {
		a.FilteredNotes = make(map[uint]models.Note, len(a.notes))
		for id, note := range a.notes {
			a.FilteredNotes[id] = note
		}
		return
	}
	query = strings.ToLower(query)
	a.FilteredNotes = make(map[uint]models.Note)
	for id, note := range a.notes {
		if strings.Contains(strings.ToLower(note.Content), query) {
			a.FilteredNotes[id] = note
			continue
		}
		for _, topic := range note.Topics {
			if strings.Contains(strings.ToLower(topic.Topic), query) {
				a.FilteredNotes[id] = note
				break
			}
		}
	}
}

// SelectCurrentNote sets the current note based on table cursor
func (a *App) SelectCurrentNote(cursor int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// Convert map to slice for cursor-based access
	notes := make([]models.Note, 0, len(a.FilteredNotes))
	for _, note := range a.FilteredNotes {
		notes = append(notes, note)
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].CreatedAt.Before(notes[j].CreatedAt)
	})
	if len(notes) == 0 || cursor >= len(notes) {
		a.currentNote = nil
		return
	}
	a.currentNote = &notes[cursor]
}

// CurrentNoteContent returns the content of the current note
func (a *App) CurrentNoteContent() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return ""
	}
	return a.currentNote.Content
}

// CurrentNoteTopics returns the topics of the current note
func (a *App) CurrentNoteTopics() []*models.Topic {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return nil
	}
	return a.currentNote.Topics
}

// HasCurrentNote checks if a note is currently selected
func (a *App) HasCurrentNote() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.currentNote != nil
}

// SaveCurrentNote updates the current note's content in-memory
func (a *App) SaveCurrentNote(content string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return
	}
	var noteID uint
	if a.currentNote.ID != 0 {
		noteID = a.currentNote.ID
	} else if tempID, exists := a.pendingNotes[a.currentNote]; exists {
		noteID = tempID
	} else {
		return // Should not happen
	}
	a.currentNote.Content = content
	a.notes[noteID] = *a.currentNote
	a.FilteredNotes[noteID] = *a.currentNote
}

// AddTopicsToCurrentNote adds topics to the current note
func (a *App) AddTopicsToCurrentNote(topicsText string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return
	}
	var noteID uint
	if a.currentNote.ID != 0 {
		noteID = a.currentNote.ID
	} else if tempID, exists := a.pendingNotes[a.currentNote]; exists {
		noteID = tempID
	} else {
		return // Should not happen
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
		}
	}
	a.notes[noteID] = *a.currentNote
	a.FilteredNotes[noteID] = *a.currentNote
}

// RemoveTopicFromCurrentNote removes a topic from the current note
func (a *App) RemoveTopicFromCurrentNote(topicToRemove string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return
	}
	var noteID uint
	if a.currentNote.ID != 0 {
		noteID = a.currentNote.ID
	} else if tempID, exists := a.pendingNotes[a.currentNote]; exists {
		noteID = tempID
	} else {
		return // Should not happen
	}

	var newTopics []*models.Topic
	for _, topic := range a.currentNote.Topics {
		if topic.Topic != topicToRemove {
			newTopics = append(newTopics, topic)
		}
	}
	a.currentNote.Topics = newTopics
	a.notes[noteID] = *a.currentNote
	a.FilteredNotes[noteID] = *a.currentNote
}

// CreateNewNote creates a new pending note
func (a *App) CreateNewNote(textareaValue string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	content := "New note"
	// if content == "" {
	// 	content = "New note"
	// }
	note := &models.Note{Content: content}
	// Assign a temporary ID for pending notes
	tempID := a.nextTempID
	a.nextTempID++
	a.pendingNotes[note] = tempID
	a.notes[tempID] = *note
	a.FilteredNotes[tempID] = *note
	a.currentNote = note
}

// DeleteCurrentNote deletes the current note in-memory
func (a *App) DeleteCurrentNote() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil || len(a.FilteredNotes) == 0 {
		return
	}
	var noteID uint
	if a.currentNote.ID != 0 {
		noteID = a.currentNote.ID
	} else if tempID, exists := a.pendingNotes[a.currentNote]; exists {
		noteID = tempID
	} else {
		return // Should not happen
	}
	// Track deletion for database sync if note was in DB
	if a.currentNote.ID != 0 && !a.isPendingNoteNoLock(a.currentNote) {
		a.deletedNotes[a.currentNote.ID] = struct{}{}
	}
	// Remove from pendingNotes if applicable
	delete(a.pendingNotes, a.currentNote)
	delete(a.notes, noteID)
	delete(a.FilteredNotes, noteID)
	a.currentNote = nil
}

// isPendingNote checks if a note is pending (public method with locking)
func (a *App) isPendingNote(note *models.Note) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.isPendingNoteNoLock(note)
}

// isPendingNoteNoLock checks if a note is pending without acquiring the mutex
func (a *App) isPendingNoteNoLock(note *models.Note) bool {
	_, exists := a.pendingNotes[note]
	return exists
}

// HasChanges checks if there are unsaved changes
func (a *App) HasChanges() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
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
	a.mutex.Lock()
	// Convert maps to slices for database sync
	notes := make([]models.Note, 0, len(a.notes))
	for _, note := range a.notes {
		notes = append(notes, note)
	}
	pendingNotes := make([]*models.Note, 0, len(a.pendingNotes))
	for note := range a.pendingNotes {
		pendingNotes = append(pendingNotes, note)
	}
	deletedNotes := make([]uint, 0, len(a.deletedNotes))
	for id := range a.deletedNotes {
		deletedNotes = append(deletedNotes, id)
	}
	a.mutex.Unlock()

	updatedNotes, err := a.db.SyncWithDatabase(notes, pendingNotes, deletedNotes)
	if err != nil {
		// Log error but continue
		return
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	// Update in-memory state
	a.notes = make(map[uint]models.Note, len(updatedNotes))
	a.FilteredNotes = make(map[uint]models.Note, len(updatedNotes))
	maxID := uint(0)
	for _, note := range updatedNotes {
		a.notes[note.ID] = note
		a.FilteredNotes[note.ID] = note
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	a.nextTempID = maxID + 1
	a.pendingNotes = make(map[*models.Note]uint)
	a.deletedNotes = make(map[uint]struct{})
}
