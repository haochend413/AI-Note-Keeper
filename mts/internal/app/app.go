package app

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/haochend413/mts/internal/db"
	"github.com/haochend413/mts/internal/models"
	"github.com/haochend413/mts/internal/types"
)

// App encapsulates application logic and state
type App struct {
	db *db.DB
	// We need both map and array, controlling the resources as pointers.
	// Notes contains everything, even the ones that are to be deleted: during sync, we first create them, and then delete them in order to keep the IDs going.
	// When filtering this, we might need to combine the list of deletednoteid to make sure the correct notes are listed.

	//////

	// New atchitecture:
	// We will keep different sorts of lists :
	// 1. all ; 2. Recent ; 3. Days ; 4. Weeks ; 5. months ; ...
	// Then we will keep a currentNote pointer of a list, switching between different lists;
	// For search, we will have a individual list that fetches from the pointer and modifies it. Basically the content of the search list depends on others.
	NotesMap            map[uint]*models.Note
	NotesList           []*models.Note // For all;
	FilteredNotesMap    map[uint]*models.Note
	FilteredNotesList   []*models.Note  // For search;
	RecentNotes         []*models.Note  // For Recent
	CurrentNotesListPtr *[]*models.Note //It points to different stuff.
	//Notes selected based on NoteSelector

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

	//init current Pointer ;
	a.CurrentNotesListPtr = &a.FilteredNotesList
	// Query the database for the maximum ID, including deleted notes
	var maxID uint
	if err := a.db.Conn.Table("notes").Select("MAX(id)").Row().Scan(&maxID); err != nil {
		log.Printf("Error fetching max ID: %v", err)
	}

	// Set the next ID for note creation
	a.nextNoteCreateID = maxID + 1
}

// This should search based on the input string, and populate the filtered list;
// The mechanism is :
// there will always be a filtered list out there, and search will switch the table content from the content of currentList to
// the filtered list.
// If not, we should populate using the NotesMap instead.
// We do not actually need the map ? maybe ? we can try it.

// This function now searches the ciurrentList and filter it to populate FilteredNotesMap and List.
func (a *App) SearchNotes(query string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	curr := *a.CurrentNotesListPtr
	if query == "" {
		// a.FilteredNotesMap = a.NotesMap
		a.FilteredNotesList = curr
		return
	}
	query = strings.ToLower(query)
	a.FilteredNotesMap = make(map[uint]*models.Note)
	a.FilteredNotesList = make([]*models.Note, 0)
	for _, note := range curr {
		if strings.Contains(strings.ToLower(note.Content), query) {
			a.FilteredNotesMap[note.ID] = note
			a.FilteredNotesList = append(a.FilteredNotesList, note)
			continue
		}
		for _, topic := range note.Topics {
			if strings.Contains(strings.ToLower(topic.Topic), query) {
				a.FilteredNotesMap[note.ID] = note
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
	if len(*a.CurrentNotesListPtr) == 0 || cursor >= len(*a.CurrentNotesListPtr) {
		a.currentNote = nil
		return
	}
	a.currentNote = (*a.CurrentNotesListPtr)[cursor]
}

// CreateNewNote creates a new pending note
func (a *App) CreateNewNote() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	note := &models.Note{Content: ""}
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	note.ID = a.nextNoteCreateID
	a.nextNoteCreateID += 1
	a.Synced = false
	a.CreateNoteIDs = append(a.CreateNoteIDs, note.ID)
	// // a.FilteredNotesMap[note.ID] = note
	a.NotesMap[note.ID] = note
	// // a.FilteredNotesList = append(a.FilteredNotesList, note)
	a.NotesList = append(a.NotesList, note)
	*a.CurrentNotesListPtr = append(*a.CurrentNotesListPtr, note)
	a.currentNote = note

}

func (a *App) UpdateRecentNotes() {
	d, _ := a.db.GetRecentNotes()
	//lock
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// Sort by ID to match display order
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
	a.RecentNotes = d
}

func (a *App) UpdateCurrentList(s types.Selector) {
	switch s {
	case types.Default:
		a.CurrentNotesListPtr = &a.NotesList
	case types.Search:
		// this should do nothing, since filtered notes list is independent of current notes list;
		// see more to this
		// but if that's the case, then many things have to be changed.
	case types.Recent:
		a.CurrentNotesListPtr = &a.RecentNotes
	}
	sort.Slice(*a.CurrentNotesListPtr, func(i, j int) bool {
		return (*a.CurrentNotesListPtr)[i].ID < (*a.CurrentNotesListPtr)[j].ID
	})
}

// SyncWithDatabase syncs in-memory changes to the database
// This only work before we sync everything.

// This is very buggy.
func (a *App) UndoDelete() {
	if len(a.DeletedNoteIDs) == 0 {
		return
	}

	lastDeletedID := a.DeletedNoteIDs[len(a.DeletedNoteIDs)-1]

	deletedNote, exists := a.NotesMap[lastDeletedID]
	if !exists {
		return
	}

	a.DeletedNoteIDs = a.DeletedNoteIDs[:len(a.DeletedNoteIDs)-1]

	a.NotesList = append(a.NotesList, deletedNote)

	if len(a.RecentNotes) < 10 {
		a.RecentNotes = append(a.RecentNotes, deletedNote)
	}
}

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

	// // Reset the filtered lists/maps to show all notes
	// a.FilteredNotesList = updatedNotes
	// a.FilteredNotesMap = make(map[uint]*models.Note, len(updatedNotes))

	// Find max ID for next note creation
	maxID := uint(0)
	for _, note := range updatedNotes {
		a.NotesMap[note.ID] = note
		// a.FilteredNotesMap[note.ID] = note
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
