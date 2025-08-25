package models

import "gorm.io/gorm"

// Note represents a note entity
type Note struct {
	gorm.Model
	Content   string
	Highlight bool
	Frequency int      //should be calculated as the number of times that is editted.
	Topics    []*Topic `gorm:"many2many:note_topics;constraint:OnDelete:CASCADE;"`
}
