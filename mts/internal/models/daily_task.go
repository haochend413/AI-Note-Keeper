package models

import "gorm.io/gorm"

// DailyTask represents a daily task entity
type DailyTask struct {
	gorm.Model
	Task    string
	Success bool
}
