package notedb

import (
	"github.com/haochend413/mantis/defs"
	"gorm.io/gorm"
)

// Clear all and then setup again; Sync database with app state
// Need to change that to accept topics
func (nd *NoteDB) SyncNoteData(notes []*defs.Note) error {
	//This might be buggy: clear table
	nd.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&defs.Note{})

	for _, n := range notes {
		if result := nd.Db.Save(n); result.Error != nil {
			return result.Error
		}
	}
	return nil
}
