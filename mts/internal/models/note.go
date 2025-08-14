package models

import "gorm.io/gorm"

// Note represents a note entity
type Note struct {
	gorm.Model
	Content string
	Topics  []*Topic `gorm:"many2many:note_topics;"`
}
