package embedder

import (
	"github.com/haochend413/ntkpr/internal/clients"
	"github.com/haochend413/ntkpr/internal/db"
	"github.com/haochend413/ntkpr/internal/models"
)

type Embedder struct {
	client *clients.EmbedClient
	db     *db.DB
}

func NewEmbedder(embclient *clients.EmbedClient, db *db.DB) *Embedder {
	return &Embedder{client: embclient, db: db}
}

func (e *Embedder) EmbedNote(noteID uint, text string) error {
	vec, err := e.client.Embed(text)
	if err != nil {
		return err
	}
	return e.db.UpsertNoteEmbedding(noteID, vec)
}

func (e *Embedder) FetchRelated(noteID uint, k int) []models.Note {
	embedding, err := e.db.GetNoteEmbedding(noteID)
	if err != nil {
		return nil
	}
	if embedding == nil {
		return nil
	}
	notes, err := e.db.SearchRelatedNotes(embedding, k)
	if err != nil {
		return nil
	}

	return notes
}
