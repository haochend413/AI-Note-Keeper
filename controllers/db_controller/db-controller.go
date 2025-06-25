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
func (m *DBManager) RefreshDaily(data []*defs.DailyTask) error {
	if err := m.DataBases.DailyDB.SyncDailyTaskData(data); err != nil {
		return err
	}
	return nil
}

func (m *DBManager) RefreshNoteTopic(data *defs.DB_Data) error {
	m.DataBases.NoteDB.SyncTopicData(data.TopicData)
	return m.DataBases.NoteDB.SyncNoteData(data.NoteData)
}

// refresh database data; Run at quit or before specific functions
func (m *DBManager) RefreshAll(data *defs.DB_Data) error {
	if err := m.DataBases.DailyDB.SyncDailyTaskData(data.DailyTaskData); err != nil {
		return err
	}
	if err := m.DataBases.NoteDB.SyncTopicData(data.TopicData); err != nil {
		return err
	}
	return m.DataBases.NoteDB.SyncNoteData(data.NoteData)
}

// fetch database data, run at the Appinit
func (m *DBManager) FetchAll() *defs.DB_Data {
	var (
		history   []defs.Note
		topics    []defs.Topic
		dailytask []defs.DailyTask
	)

	// Fetch notes and daily tasks, handle errors
	if err := m.DataBases.NoteDB.Db.Find(&history).Error; err != nil {
		return &defs.DB_Data{NoteData: []*defs.Note{}}
	}
	if err := m.DataBases.NoteDB.Db.Find(&topics).Error; err != nil {
		return &defs.DB_Data{TopicData: []*defs.Topic{}}
	}
	if err := m.DataBases.DailyDB.Db.Find(&dailytask).Error; err != nil {
		return &defs.DB_Data{DailyTaskData: []*defs.DailyTask{}}
	}

	//value-pointer conversion

	notePtrs := make([]*defs.Note, 0, len(history))
	dailytaskPtrs := make([]*defs.DailyTask, 0, len(dailytask))
	topicPtrs := make([]*defs.Topic, 0, len(topics))
	if len(history) != 0 {
		for i := range history {
			notePtrs = append(notePtrs, &history[i])
		}
	}
	if len(dailytask) != 0 {
		for i := range dailytask {
			dailytaskPtrs = append(dailytaskPtrs, &dailytask[i])
		}
	}
	if len(topics) != 0 {
		for i := range topics {
			topicPtrs = append(topicPtrs, &topics[i])
		}
	}

	return &defs.DB_Data{NoteData: notePtrs, DailyTaskData: dailytaskPtrs, TopicData: topicPtrs}
}
