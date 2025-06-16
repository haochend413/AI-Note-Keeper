package notedb

import (
	"gorm.io/gorm"
)

type NoteDB struct {
	Db *gorm.DB
	// Name string
}

// func (nd *NoteDB) Init() error {
// 	NoteDB.Db
// }

func (nd *NoteDB) Close() error {
	n, err := nd.Db.DB()
	if err != nil {
		return err
	}
	return n.Close()
}

func (nd *NoteDB) GetDB() *gorm.DB {
	return nd.Db
}
