package db

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

// struct for single message
type Note struct {
	gorm.Model
	Content string
}
