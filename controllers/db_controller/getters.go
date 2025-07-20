package dbcontroller

import (
	"github.com/haochend413/mantis/defs"
)

// This should be fast;
func (m *DBManager) FetchNoteFromID(id int) (*defs.Note, error) {
	if id == 0 {
		return nil, nil
	}
	var note defs.Note
	result := m.DataBases.NoteDB.Db.Preload("Topics").First(&note, id)
	return &note, result.Error
}
