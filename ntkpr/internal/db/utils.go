package db

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/haochend413/ntkpr/internal/models"
)

type NextCreateIDs struct {
	NoteID   uint
	BranchID uint
	ThreadID uint
}

func (d *DB) GetNextCreateIDs() NextCreateIDs {
	var ids NextCreateIDs
	var noteMax, branchMax, threadMax uint
	d.Conn.Raw("SELECT COALESCE(MAX(id), 0) FROM notes").Scan(&noteMax)
	d.Conn.Raw("SELECT COALESCE(MAX(id), 0) FROM branches").Scan(&branchMax)
	d.Conn.Raw("SELECT COALESCE(MAX(id), 0) FROM threads").Scan(&threadMax)
	ids.NoteID = noteMax + 1
	ids.BranchID = branchMax + 1
	ids.ThreadID = threadMax + 1
	return ids
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
