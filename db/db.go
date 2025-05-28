package db

import (
	"gorm.io/gorm"
)

type NoteDB struct {
	Db *gorm.DB
}

// var _ db_models.DBController = (*NoteDB)(nil)

// func (nd *NoteDB) DBInit() {
// 	// open notes database
// 	var err error
// 	n, err := gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	//assign
// 	nd.Db = n
// 	n.AutoMigrate(&Note{})
// }

// func (nd *NoteDB) DBAdd(content string) error {
// 	//init note struct
// 	note := &Note{Content: content}
// 	//pass the string to database;
// 	result := nd.Db.Create(note)
// 	return result.Error
// }
