package db

import (
	"log"
	"strings"

	"github.com/haochend413/mts/internal/models"
)

// SyncWithDatabase persists in-memory changes to the database
func (d *DB) SyncWithDatabase(notes []models.Note, pendingNotes []*models.Note, deletedNotes []uint) ([]models.Note, error) {
	// Delete notes marked for deletion
	for _, noteID := range deletedNotes {
		d.Conn.Delete(&models.Note{}, noteID)
	}

	// Save pending notes
	for _, note := range pendingNotes {
		note.Content = strings.TrimSpace(note.Content)
		if note.Content == "" {
			note.Content = "New note"
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
	} // Load all notes from database
	var dbNotes []models.Note
	if err := d.Conn.Preload("Topics").Find(&dbNotes).Error; err != nil {
		return nil, err
	}
	return dbNotes, nil
}
