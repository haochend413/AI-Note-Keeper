package models

import "gorm.io/gorm"

// Topic represents a topic entity
type Topic struct {
	gorm.Model
	Topic string
	Notes []*Note `gorm:"many2many:note_topics;"`
}
