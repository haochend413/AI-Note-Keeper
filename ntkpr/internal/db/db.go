package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/haochend413/ntkpr/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps the GORM database connection
type DB struct {
	Conn *gorm.DB
}

// NewDB initializes a new database connection and migrates schema
func NewDB(path string) (*DB, error) {
	// if not exist, create all dirs
	_, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, create directory and config file with defaults
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating database directory: %v", err)
			return nil, err
		}

	}

	conn, err := gorm.Open(sqlite.Open(path+"?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate schema
	err = conn.AutoMigrate(&models.Note{}, &models.Thread{}, &models.Branch{})
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
