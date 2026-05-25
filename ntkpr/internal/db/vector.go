package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/haochend413/ntkpr/internal/models"
)

const EmbeddingDim = 1024

type RelatedNote struct {
	Note     models.Note
	Distance float64
}

func (d *DB) rawDB() (*sql.DB, error) {
	return d.Conn.DB()
}

func (d *DB) InitVectorTable() error {
	sqlDB, err := d.rawDB()
	if err != nil {
		return err
	}

	var version string
	if err := sqlDB.QueryRow(`SELECT vec_version()`).Scan(&version); err != nil {
		return fmt.Errorf("sqlite-vec not loaded: %w", err)
	}

	_, err = sqlDB.Exec(`
CREATE VIRTUAL TABLE IF NOT EXISTS note_vecs USING vec0(
    embedding FLOAT[1024]
);
`)
	return err
}

func (d *DB) UpsertNoteEmbedding(noteID uint, embedding []float32) error {
	if len(embedding) != EmbeddingDim {
		return fmt.Errorf("embedding dim mismatch: got %d, want %d", len(embedding), EmbeddingDim)
	}

	b, err := json.Marshal(embedding)
	if err != nil {
		return err
	}

	sqlDB, err := d.rawDB()
	if err != nil {
		return err
	}

	_, err = sqlDB.Exec(`
INSERT OR REPLACE INTO note_vecs(rowid, embedding)
VALUES (?, ?)
`, noteID, string(b))

	return err
}

func (d *DB) SearchRelatedNotes(queryEmbedding []float32, k int) ([]RelatedNote, error) {
	if len(queryEmbedding) != EmbeddingDim {
		return nil, fmt.Errorf("embedding dim mismatch: got %d, want %d", len(queryEmbedding), EmbeddingDim)
	}

	b, err := json.Marshal(queryEmbedding)
	if err != nil {
		return nil, err
	}

	sqlDB, err := d.rawDB()
	if err != nil {
		return nil, err
	}

	rows, err := sqlDB.Query(`
SELECT rowid, distance
FROM note_vecs
WHERE embedding MATCH ?
ORDER BY distance
LIMIT ?
`, string(b), k)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []RelatedNote

	for rows.Next() {
		var noteID uint
		var distance float64

		if err := rows.Scan(&noteID, &distance); err != nil {
			return nil, err
		}

		var note models.Note
		if err := d.Conn.Preload("Branches").First(&note, noteID).Error; err != nil {
			return nil, err
		}

		results = append(results, RelatedNote{
			Note:     note,
			Distance: distance,
		})
	}

	return results, rows.Err()
}
