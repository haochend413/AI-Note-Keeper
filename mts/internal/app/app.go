package app

import (
	"log"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/haochend413/mts/internal/db"
	"github.com/haochend413/mts/internal/models"
)

// App encapsulates application logic and state
type App struct {
	db *db.DB
	// We need both map and array, controlling the resources as pointers.
	// Notes contains everything, even the ones that are to be deleted: during sync, we first create them, and then delete them in order to keep the IDs going.
	// filtered notes are all the notes are notes that are shown on the screen.
	// When filtering this, we might need to combine the list of deletednoteid to make sure the correct notes are listed.
	NotesMap          map[uint]*models.Note
	NotesList         []*models.Note
	FilteredNotesMap  map[uint]*models.Note
	FilteredNotesList []*models.Note
	RecentNotes       []*models.Note
	// The current note that is selected, in order to change and demo;
	currentNote *models.Note
	// In order to manage pending notes, We record the changed note IDs, and we send them back to database;
	PendingNoteIDs []uint // ok on create we
	DeletedNoteIDs []uint //
	CreateNoteIDs  []uint

	// This should be one larger than the last note i have in my db;
	nextNoteCreateID uint
	Synced           bool
	// Topics
	Topics map[uint]*models.Topic
	mutex  sync.Mutex
}

// NewApp creates a new application instance
func NewApp(dbConn *db.DB) *App {
	app := &App{
		db:                dbConn,
		NotesMap:          make(map[uint]*models.Note),
		NotesList:         make([]*models.Note, 0),
		FilteredNotesMap:  make(map[uint]*models.Note),
		FilteredNotesList: make([]*models.Note, 0),
		RecentNotes:       make([]*models.Note, 0),
		Topics:            make(map[uint]*models.Topic),
		nextNoteCreateID:  1, // Default starting ID if database is empty
		PendingNoteIDs:    []uint{},
		DeletedNoteIDs:    []uint{},
	}

	//load everything into the app
	app.loadData()
	// print(len(app.NotesList))
	// print(len(app.NotesMap))
	// print(len(app.FilteredNotesList))
	// print(len(app.FilteredNotesMap))
	return app
}

// loadNotes loads notes from the database
// This function also load all topics
func (a *App) loadData() {
	//fetch data from database;
	notes, topics, err := a.db.SyncData(make(map[uint]*models.Note), nil, nil, nil)
	if err != nil {
		// Log error but continue with empty notes
		log.Panic(err)
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// fill in the variables
	a.FilteredNotesList = notes
	a.NotesList = notes
	a.NotesMap = make(map[uint]*models.Note, len(notes))
	a.Topics = make(map[uint]*models.Topic, len(topics))
	a.FilteredNotesMap = make(map[uint]*models.Note, len(notes))
	// init recent
	recentNotes, _ := a.db.GetRecentNotes()
	a.RecentNotes = recentNotes
	for _, note := range notes {
		a.NotesMap[note.ID] = note
		a.FilteredNotesMap[note.ID] = note
	}
	for _, topic := range topics {
		a.Topics[topic.ID] = topic
	}
	// Query the database for the maximum ID, including deleted notes
	var maxID uint
	if err := a.db.Conn.Table("notes").Select("MAX(id)").Row().Scan(&maxID); err != nil {
		log.Printf("Error fetching max ID: %v", err)
	}

	// Set the next ID for note creation
	a.nextNoteCreateID = maxID + 1
}

// SearchNotes filters notes based on a query
// This should remain the same; s
// SearchNotes filters notes based on a query
func (a *App) SearchNotes(query string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if query == "" {
		a.FilteredNotesMap = a.NotesMap
		a.FilteredNotesList = a.NotesList
		return
	}
	query = strings.ToLower(query)
	a.FilteredNotesMap = make(map[uint]*models.Note)
	a.FilteredNotesList = make([]*models.Note, 0)
	for id, note := range a.NotesMap {
		if strings.Contains(strings.ToLower(note.Content), query) {
			a.FilteredNotesMap[id] = note
			a.FilteredNotesList = append(a.FilteredNotesList, note)
			continue
		}
		for _, topic := range note.Topics {
			if strings.Contains(strings.ToLower(topic.Topic), query) {
				a.FilteredNotesMap[id] = note
				a.FilteredNotesList = append(a.FilteredNotesList, note)
				break
			}
		}
	}
	sort.Slice(a.FilteredNotesList, func(i, j int) bool {
		return a.FilteredNotesList[i].ID < a.FilteredNotesList[j].ID
	})
}

// SelectCurrentNote sets the current note based on table cursor
// This one should be slow since it turns. Maybe consider using more space. We should do this at the start of the program.
func (a *App) SelectCurrentNote(cursor int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	//This might be wrong: 1. needs sort;
	if len(a.FilteredNotesList) == 0 || cursor >= len(a.FilteredNotesList) {
		a.currentNote = nil
		return
	}
	a.currentNote = a.FilteredNotesList[cursor]
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

// CurrentNoteTopics returns the topics of the current note
func (a *App) CurrentNoteID() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return -1
	}
	return int(a.currentNote.ID)
}

func (a *App) CurrentNoteLastUpdate() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return time.Time{}
	}
	return a.currentNote.UpdatedAt
}

func (a *App) CurrentNoteFrequency() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return 0
	}
	return a.currentNote.Frequency
}

// HasCurrentNote checks if a note is currently selected
func (a *App) HasCurrentNote() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.currentNote != nil
}

// Update the content of current note, content fetched from terminal
// No need to update the
func (a *App) SaveCurrentNote(content string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.currentNote == nil {
		return
	}
	var noteID uint
	if a.currentNote.ID != 0 {
		noteID = a.currentNote.ID
	}
	if a.currentNote.Content != content {
		a.currentNote.Content = content
		a.currentNote.Frequency += 1
		a.Synced = false
		a.PendingNoteIDs = append(a.PendingNoteIDs, noteID)
	}
	// a.NotesMap[noteID] = a.currentNote

	// a.FilteredNotesMap[noteID] = a.currentNote

	//push into recent edited notes;
}

func (a *App) HighlightCurrentNote() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.HasCurrentNote() {
		a.currentNote.Highlight = !a.currentNote.Highlight
	}
}

// AddTopicsToCurrentNote adds topics to the current note
func (a *App) AddTopicsToCurrentNote(topicsText string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
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
		}

	}
	//mark as pending
	a.Synced = false

	a.PendingNoteIDs = append(a.PendingNoteIDs, a.currentNote.ID)
	// a.notes[noteID] = *a.currentNote
	// a.FilteredNotes[noteID] = *a.currentNote
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
	}

	var newTopics []*models.Topic
	for _, topic := range a.currentNote.Topics {
		if topic.Topic != topicToRemove {
			newTopics = append(newTopics, topic)
		}
	}
	a.currentNote.Topics = newTopics
	a.Synced = false

	a.PendingNoteIDs = append(a.PendingNoteIDs, noteID)
}

// CreateNewNote creates a new pending note
func (a *App) CreateNewNote() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	note := &models.Note{Content: ""}
	note.ID = a.nextNoteCreateID
	a.nextNoteCreateID += 1
	a.Synced = false

	a.CreateNoteIDs = append(a.CreateNoteIDs, note.ID)
	a.FilteredNotesMap[note.ID] = note
	a.NotesMap[note.ID] = note
	a.FilteredNotesList = append(a.FilteredNotesList, note)
	a.NotesList = append(a.NotesList, note)
	a.currentNote = note

}

// DeleteCurrentNote deletes the current note in-memory
func (a *App) DeleteCurrentNote(cursor uint) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Check if there's a current note
	if a.currentNote == nil {
		return
	}

	// Get the note ID
	noteID := a.currentNote.ID

	// Handle differently based on note status
	isInCreateList := slices.Contains(a.CreateNoteIDs, noteID)

	if isInCreateList {
		for i, id := range a.CreateNoteIDs {
			if id == noteID {
				a.CreateNoteIDs = append(a.CreateNoteIDs[:i], a.CreateNoteIDs[i+1:]...)
				//Remove it from the FilteredNotesList and Notes List
				break
			}
		}
	} else if noteID != 0 {
		a.DeletedNoteIDs = append(a.DeletedNoteIDs, noteID)
		// Also remove from pending if it was pending
		// for i, id := range a.PendingNoteIDs {
		// 	if id == noteID {
		// 		a.PendingNoteIDs = append(a.PendingNoteIDs[:i], a.PendingNoteIDs[i+1:]...)
		// 		break
		// 	}
		// }
	}

	// // Remove from NotesMap
	// delete(a.NotesMap, noteID)

	// // Remove from NotesList
	// for i, note := range a.NotesList {
	// 	if note.ID == noteID {
	// 		a.NotesList = append(a.NotesList[:i], a.NotesList[i+1:]...)
	// 		break
	// 	}
	// }

	delete(a.FilteredNotesMap, noteID)
	for i, note := range a.FilteredNotesList {
		if note.ID == noteID {
			a.FilteredNotesList = append(a.FilteredNotesList[:i], a.FilteredNotesList[i+1:]...)
			break
		}
	}
	a.Synced = false

	// Clear the current note reference
	// This might need debugging and border conditions management;
	if cursor >= uint(len(a.FilteredNotesList)) {
		cursor = uint(len(a.FilteredNotesList) - 1)
	}
	a.currentNote = a.FilteredNotesList[cursor]
}

func (a *App) UpdateRecentNotes() {
	d, _ := a.db.GetRecentNotes()
	//lock
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.RecentNotes = d
}

// SyncWithDatabase syncs in-memory changes to the database

func (a *App) SyncWithDatabase() {
	a.mutex.Lock()
	// Make a copy of the IDs slices for the database operation
	pendingIDs := append([]uint{}, a.PendingNoteIDs...)
	deletedIDs := append([]uint{}, a.DeletedNoteIDs...)
	createIDs := append([]uint{}, a.CreateNoteIDs...)

	// Make a copy of the notes map
	notesMapCopy := make(map[uint]*models.Note, len(a.NotesMap))
	for id, note := range a.NotesMap {
		notesMapCopy[id] = note
	}
	a.mutex.Unlock()

	// Sync with the database
	updatedNotes, updatedTopics, err := a.db.SyncData(notesMapCopy, pendingIDs, deletedIDs, createIDs)
	if err != nil {
		log.Printf("Error syncing with database: %v", err)
		return
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Update in-memory state with the fresh data from database
	a.NotesList = updatedNotes
	a.NotesMap = make(map[uint]*models.Note, len(updatedNotes))
	a.Topics = make(map[uint]*models.Topic, len(updatedTopics))

	// Reset the filtered lists/maps to show all notes
	a.FilteredNotesList = updatedNotes
	a.FilteredNotesMap = make(map[uint]*models.Note, len(updatedNotes))

	// Find max ID for next note creation
	maxID := uint(0)
	for _, note := range updatedNotes {
		a.NotesMap[note.ID] = note
		a.FilteredNotesMap[note.ID] = note
		if note.ID > maxID {
			maxID = note.ID
		}
	}

	// Update topics
	for _, topic := range updatedTopics {
		a.Topics[topic.ID] = topic
	}

	// Set next ID and clear pending/deleted lists
	a.nextNoteCreateID = maxID + 1
	a.PendingNoteIDs = []uint{}
	a.DeletedNoteIDs = []uint{}
	a.CreateNoteIDs = []uint{}
	a.Synced = true

	// If we had a current note, try to find it in the updated notes
	if a.currentNote != nil {
		currentID := a.currentNote.ID
		if note, exists := a.NotesMap[currentID]; exists {
			a.currentNote = note
		} else if len(a.NotesList) > 0 {
			// If current note was deleted, select the first note
			a.currentNote = a.NotesList[0]
		} else {
			a.currentNote = nil
		}
	}
}
