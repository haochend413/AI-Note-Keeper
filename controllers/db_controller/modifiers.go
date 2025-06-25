package dbcontroller

import "github.com/haochend413/mantis/defs"

// This is a simple version:
// No lazy-sync, directly interact with db.
// For cmd that ends after one instruction.
// No possibility of in-memory unstored notes / topics.
// Need help: a way to view all the notes and topics with their ids.
func (m *DBManager) LinkNoteTopic(noteid string, topicid string) error {
	var note defs.Note
	if err := m.DataBases.NoteDB.Db.First(&note, noteid).Error; err != nil {
		return err
	}
	var topic defs.Topic
	if err := m.DataBases.NoteDB.Db.First(&topic, topicid).Error; err != nil {
		return err
	}
	if err := m.DataBases.NoteDB.Db.Model(&note).Association("Topics").Append(&topic); err != nil {
		return err
	}
	return nil
}
