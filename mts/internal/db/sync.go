package db

import (
	"log"
	"strings"

	"github.com/haochend413/mts/internal/models"
)

// SyncWithDatabase persists in-memory changes to the database
func (d *DB) SyncWithDatabase(notes []models.Note, pendingNotes []*models.Note, deletedNotes []uint) ([]models.Note, []models.Topic, error) {
	// Delete notes marked for deletion
	for _, noteID := range deletedNotes {
		d.Conn.Delete(&models.Note{}, noteID)
	}

	// Save pending notes
	for _, note := range pendingNotes {
		note.Content = strings.TrimSpace(note.Content)
		if note.Content == "" {
			note.Content = ""
		}

		// First create the note to get its ID
		result := d.Conn.Create(note)
		if result.Error != nil {
			log.Printf("Error creating note in DB: %v", result.Error)
			continue
		}

		// Handle topics separately after creation
		if len(note.Topics) > 0 {
			// Reload the note to ensure we have the correct associations
			var createdNote models.Note
			if err := d.Conn.First(&createdNote, note.ID).Error; err != nil {
				log.Printf("Error finding created note: %v", err)
				continue
			}

			// Now handle the topics
			d.Conn.Model(&createdNote).Association("Topics").Clear()
			d.Conn.Model(&createdNote).Association("Topics").Append(note.Topics)
		}
	}

	// Save updated notes
	for i := range notes {
		if notes[i].ID != 0 { // Only save notes that were previously in DB
			// First find the existing note to properly update it
			var existingNote models.Note
			if err := d.Conn.First(&existingNote, notes[i].ID).Error; err != nil {
				log.Printf("Error finding note with ID %d: %v", notes[i].ID, err)
				continue
			}
			// Update the note's content
			existingNote.Content = notes[i].Content
			d.Conn.Save(&existingNote)

			// Handle the topics association
			d.Conn.Model(&existingNote).Association("Topics").Clear()
			if len(notes[i].Topics) > 0 {
				d.Conn.Model(&existingNote).Association("Topics").Append(notes[i].Topics)
			}
		}
	}

	// Load all notes from database
	var dbNotes []models.Note
	if err := d.Conn.Preload("Topics").Find(&dbNotes).Error; err != nil {
		return nil, nil, err
	}
	var dbTopics []models.Topic
	if err := d.Conn.Preload("Notes").Find(&dbTopics).Error; err != nil {
		return nil, nil, err
	}
	return dbNotes, dbTopics, nil
}
func (d *DB) SyncData(NotesMap map[uint]*models.Note, PendingNoteIDs []uint, DeletedNoteIDs []uint, CreateNoteIDs []uint) ([]*models.Note, []*models.Topic, error) {
	// Combine create and pending notes - they're handled identically
	notesToSave := append([]uint{}, CreateNoteIDs...)
	notesToSave = append(notesToSave, PendingNoteIDs...)

	// Process all notes to save
	for _, id := range notesToSave {
		note := NotesMap[id]
		if note == nil {
			continue
		}

		note.Content = strings.TrimSpace(note.Content)
		if note.Content == "" {
			note.Content = ""
		}

		// Save the note
		result := d.Conn.Save(note)
		if result.Error != nil {
			log.Printf("Error creating note in DB: %v", result.Error)
			continue
		}

		// Ensure all topics exist in the database
		for i, topic := range note.Topics {
			if topic.ID == 0 {
				// Check if the topic already exists in the database
				var existingTopic models.Topic
				if err := d.Conn.Where("topic = ?", topic.Topic).First(&existingTopic).Error; err == nil {
					// Use the existing topic
					note.Topics[i] = &existingTopic
				} else {
					// Create a new topic
					if err := d.Conn.Create(topic).Error; err != nil {
						log.Printf("Error creating topic '%s': %v", topic.Topic, err)
						continue
					}
					log.Printf("Created new topic: %s, ID: %d", topic.Topic, topic.ID)
				}
			}
		}

		// Reload the note to ensure we have the correct associations
		var createdNote models.Note
		if err := d.Conn.First(&createdNote, note.ID).Error; err != nil {
			log.Printf("Error finding created note: %v", err)
			continue
		}

		// Associate topics with the note
		if err := d.Conn.Model(&createdNote).Association("Topics").Clear(); err != nil {
			log.Printf("Error clearing topics for note ID %d: %v", createdNote.ID, err)
			continue
		}
		if err := d.Conn.Model(&createdNote).Association("Topics").Append(note.Topics); err != nil {
			log.Printf("Error associating topics with note ID %d: %v", createdNote.ID, err)
			continue
		}
	}

	// Delete notes marked for deletion
	for _, noteID := range DeletedNoteIDs {
		if err := d.Conn.Delete(&models.Note{}, noteID).Error; err != nil {
			return nil, nil, err
		}
	}

	// Load all notes from database
	var dbNotes []*models.Note
	if err := d.Conn.Preload("Topics").Find(&dbNotes).Error; err != nil {
		return nil, nil, err
	}

	// Load all topics from database
	var dbTopics []*models.Topic
	if err := d.Conn.Preload("Notes").Find(&dbTopics).Error; err != nil {
		return nil, nil, err
	}

	return dbNotes, dbTopics, nil
}
