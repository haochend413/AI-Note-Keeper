package db

import (
	"github.com/haochend413/mts/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps the GORM database connection
type DB struct {
	Conn *gorm.DB
}

// NewDB initializes a new database connection and migrates schema
func NewDB(path string) (*DB, error) {
	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate schema
	err = conn.AutoMigrate(&models.Note{}, &models.Topic{}, &models.DailyTask{})
	if err != nil {
		return nil, err
	}
	return &DB{Conn: conn}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	sqlDB, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
