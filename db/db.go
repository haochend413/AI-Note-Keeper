package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var NoteDB *gorm.DB

func NoteDBInit() {
	// open notes database
	var err error
	NoteDB, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}
	NoteDB.AutoMigrate(&Note{})
}
