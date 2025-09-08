package app

import (
	"slices"
	"strings"
	"time"

	"github.com/haochend413/mts/internal/models"
)

//is the mutex really required ? Well, maybe making current note public is a good idea, this is just stupid.

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

	topicsText = strings.ToLower(strings.TrimSpace(topicsText))
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

	// so....it is still in notesMap;

	// delete(*a.CurrentNotesListPtr, noteID)
	for i, note := range *a.CurrentNotesListPtr {
		if note.ID == noteID {
			*a.CurrentNotesListPtr = append((*a.CurrentNotesListPtr)[:i], (*a.CurrentNotesListPtr)[i+1:]...)
			break
		}
	}
	a.Synced = false

	// Clear the current note reference
	// This might need debugging and border conditions management;
	if cursor >= uint(len(*a.CurrentNotesListPtr)) {
		cursor = uint(len(*a.CurrentNotesListPtr) - 1)
	}
	a.currentNote = (*a.CurrentNotesListPtr)[cursor]
}
