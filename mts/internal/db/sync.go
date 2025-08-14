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
		result := d.Conn.Create(note)
		if result.Error != nil {
			log.Printf("Error creating note in DB: %v", result.Error)
			continue
		}
		if len(note.Topics) > 0 {
			d.Conn.Model(note).Association("Topics").Append(note.Topics)
		}
	}

	// Save updated notes
	for _, note := range notes {
		if note.ID != 0 { // Only save notes that were previously in DB
			d.Conn.Save(&note)
		}
	}

	// Load all notes from database
	var dbNotes []models.Note
	if err := d.Conn.Preload("Topics").Find(&dbNotes).Error; err != nil {
		return nil, err
	}
	return dbNotes, nil
}
