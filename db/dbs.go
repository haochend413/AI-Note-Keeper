package db

import (
	"github.com/haochend413/mantis/db/notedb"
)

// var DBs *DataBases

type DataBases struct {
	NoteDB *notedb.NoteDB
}

func (DBs *DataBases) InitAll() {
	DBs.NoteDB = &notedb.NoteDB{}
	DBs.NoteDB.Db = DBInit("notes")
}

func (DBs *DataBases) CloseAll() {
	_ = DBs.NoteDB.Close()
}
