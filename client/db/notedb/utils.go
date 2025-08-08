package notedb

import (
	"github.com/haochend413/mantis/defs"
	"gorm.io/gorm"
)

// Clear all and then setup again; Sync database with app state
// Need to change that to accept topics
func (nd *NoteDB) SyncNoteData(notes []*defs.Note) error {
	nd.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&defs.Note{})
	for _, n := range notes {
		if result := nd.Db.Save(n); result.Error != nil {
			return result.Error
		}
		// Sync the many-to-many relationship
		if err := nd.Db.Model(n).Association("Topics").Replace(n.Topics); err != nil {
			return err
		}
	}
	return nil
}

func (nd *NoteDB) SyncTopicData(topics []*defs.Topic) error {
	nd.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&defs.Topic{})
	for _, t := range topics {
		if result := nd.Db.Save(t); result.Error != nil {
			return result.Error
		}
		// Sync the many-to-many relationship
		if err := nd.Db.Model(t).Association("Notes").Replace(t.Notes); err != nil {
			return err
		}
	}
	return nil
}
