package db

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"

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

	if _, err = sqlDB.Exec(`DELETE FROM note_vecs WHERE rowid = ?`, noteID); err != nil {
		return err
	}
	_, err = sqlDB.Exec(`INSERT INTO note_vecs(rowid, embedding) VALUES (?, ?)`, noteID, string(b))
	return err
}

// GetNoteEmbedding fetches the stored embedding for a note. Returns nil with no error if the note
// has no embedding yet (not yet synced to the vector table).
func (d *DB) GetNoteEmbedding(noteID uint) ([]float32, error) {
	sqlDB, err := d.rawDB()
	if err != nil {
		return nil, err
	}

	var blob []byte
	err = sqlDB.QueryRow(`SELECT embedding FROM note_vecs WHERE rowid = ?`, noteID).Scan(&blob)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(blob)%4 != 0 {
		return nil, fmt.Errorf("unexpected embedding blob length: %d bytes", len(blob))
	}
	embedding := make([]float32, len(blob)/4)
	for i := range embedding {
		bits := binary.LittleEndian.Uint32(blob[i*4 : i*4+4])
		embedding[i] = math.Float32frombits(bits)
	}
	return embedding, nil
}

func (d *DB) SearchRelatedNotes(queryEmbedding []float32, k int) ([]models.Note, error) {
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

	var results []models.Note

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

		results = append(results, note)
	}

	return results, rows.Err()
}
