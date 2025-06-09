//This package should import db and design functions, generate output for gui to render

//Lazy sync & writeback structure: controller should be able to write data from db to the local copy in main; (all sorts after managed by DB operations);
//And it should also be able to send gui data back to db to re-fresh the states;

package dbcontroller

import (
	"github.com/haochend413/mantis/db"
	"github.com/haochend413/mantis/models"
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
func (m *DBManager) RefreshAll(data *models.DB_Data) error {
	return m.DataBases.NoteDB.SyncNoteData(data.NoteDBData)
}

// fetch database data, run at the Appinit
func (m *DBManager) FetchAll() *models.DB_Data {
	var history []models.Note
	result := m.DataBases.NoteDB.Db.Find(&history)
	if result.Error != nil {
		// handle error properly (optional)
		return &models.DB_Data{NoteDBData: []*models.Note{}}
	}

	//value-pointer conversion
	notePtrs := make([]*models.Note, 0, len(history))
	for i := range history {
		notePtrs = append(notePtrs, &history[i])
	}

	return &models.DB_Data{NoteDBData: notePtrs}
}
