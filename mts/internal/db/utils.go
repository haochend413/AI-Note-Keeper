package db

import "github.com/haochend413/mts/internal/models"

func (d *DB) GetRecentNotes() ([]*models.Note, error) {
	var recentNotes []*models.Note
	// I do not need their topics.
	if err := d.Conn.Preload("Topics").Order("ID DESC").Limit(10).Find(&recentNotes).Error; err != nil {
		return nil, err
	}
	return recentNotes, nil
}
