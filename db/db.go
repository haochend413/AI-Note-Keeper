package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var noteDB *gorm.DB

func NoteDBInit() {
	// open notes database
	var err error
	noteDB, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}
	noteDB.AutoMigrate(&Note{})
}
