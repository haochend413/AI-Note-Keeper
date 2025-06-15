package models

import (
	"gorm.io/gorm"
)

// // gorm.Model definition
// type Model struct {
//   ID        uint           `gorm:"primaryKey"`
//   CreatedAt time.Time
//   UpdatedAt time.Time
//   DeletedAt gorm.DeletedAt `gorm:"index"`
// }

// structure for db data storage
type DB_Data struct {
	NoteData      []*Note
	TopicData     []*Topic
	NoteTopicData []*NoteTopic
}

// struct for single message
type Note struct {
	// use the unique ID as indicator
	gorm.Model
	Content string
	Topics  []*Topic `gorm:"many2many:note_topics;"`
}

type Topic struct {
	gorm.Model
	Topic string
	Notes []*Note `gorm:"many2many:note_topics;"`
}

// link between noteID and topicID
type NoteTopic struct {
	NoteID  uint
	TopicID uint
}
