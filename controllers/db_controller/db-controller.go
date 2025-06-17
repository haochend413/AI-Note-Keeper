//This package should import db and design functions, generate output for gui to render

//Lazy sync & writeback structure: controller should be able to write data from db to the local copy in main; (all sorts after managed by DB operations);
//And it should also be able to send gui data back to db to re-fresh the states;

package dbcontroller

import (
	"github.com/haochend413/mantis/db"
	"github.com/haochend413/mantis/defs"
)

type DBManager struct {
	DataBases *db.DataBases
	// Controller should get and set the temporary data stored inside the gui component;
	NeedRefresh bool
}

func (m *DBManager) InitManager() error {
	m.DataBases = &db.DataBases{}
	m.DataBases.InitAll()
	return nil
}

func (m *DBManager) CloseManager() error {
	m.DataBases.CloseAll()
	return nil
}

// refresh database data; Run at quit or before specific functions
func (m *DBManager) RefreshAll(data *defs.DB_Data) error {
	return m.DataBases.NoteDB.SyncNoteData(data.NoteData)
}

// fetch database data, run at the Appinit
func (m *DBManager) FetchAll() *defs.DB_Data {
	var history []defs.Note
	result := m.DataBases.NoteDB.Db.Find(&history)
	if result.Error != nil {
		// handle error properly (optional)
		return &defs.DB_Data{NoteData: []*defs.Note{}}
	}

	//value-pointer conversion
	notePtrs := make([]*defs.Note, 0, len(history))
	for i := range history {
		notePtrs = append(notePtrs, &history[i])
	}

	return &defs.DB_Data{NoteData: notePtrs}
}
