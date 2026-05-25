package db

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/haochend413/ntkpr/internal/models"
)

func (d *DB) GetCreateNoteID() uint {
	// Query the database for the maximum ID, including deleted notes
	var maxID uint
	if err := d.Conn.Table("notes").Select("MAX(id)").Row().Scan(&maxID); err != nil {
		maxID = 0
	}
	return maxID + 1
}

func (d *DB) GetCreateBranchID() uint {
	// Query the database for the maximum ID, including deleted notes
	var maxID uint
	if err := d.Conn.Table("branches").Select("MAX(id)").Row().Scan(&maxID); err != nil {
		maxID = 0
	}
	return maxID + 1
}
func (d *DB) GetCreateThreadID() uint {
	// Query the database for the maximum ID, including deleted notes
	var maxID uint
	if err := d.Conn.Table("threads").Select("MAX(id)").Row().Scan(&maxID); err != nil {
		maxID = 0
	}
	return maxID + 1
}

// export the serialized data into desired position
func (d *DB) ExportNoteToJSON(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	var notes []models.Note
	err := d.Conn.Find(&notes).Error
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
