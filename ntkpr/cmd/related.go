package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/haochend413/ntkpr/internal/models"
	"github.com/spf13/cobra"
)

func resetAndReembedAllNotes() error {
	if globalEmbedClient == nil {
		return fmt.Errorf("no embedding client configured")
	}

	sqlDB, err := globalDB.Conn.DB()
	if err != nil {
		return err
	}

	fmt.Println("Dropping old note_vecs...")
	if _, err := sqlDB.Exec(`DROP TABLE IF EXISTS note_vecs;`); err != nil {
		return err
	}

	fmt.Println("Creating note_vecs...")
	if _, err := sqlDB.Exec(`
CREATE VIRTUAL TABLE note_vecs USING vec0(
	embedding FLOAT[1024]
);
`); err != nil {
		return err
	}

	var notes []models.Note
	if err := globalDB.Conn.Find(&notes).Error; err != nil {
		return err
	}

	fmt.Printf("Re-embedding %d notes...\n", len(notes))

	count := 0
	for _, note := range notes {
		content := strings.TrimSpace(note.Content)
		if content == "" {
			continue
		}

		vec, err := globalEmbedClient.Embed(content)
		if err != nil {
			return fmt.Errorf("embed failed for note %d: %w", note.ID, err)
		}

		if err := globalDB.UpsertNoteEmbedding(note.ID, vec); err != nil {
			return fmt.Errorf("upsert failed for note %d: %w", note.ID, err)
		}

		count++
		if count%20 == 0 {
			fmt.Printf("  embedded %d notes...\n", count)
		}
	}

	fmt.Printf("Done. Re-embedded %d non-empty notes.\n", count)
	return nil
}

var relatedResetReembed bool

var RelatedNotesCmd = &cobra.Command{
	Use:   "related <note-id>",
	Short: "Print a note and its top 5 semantically related notes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid note ID %q: %v\n", args[0], err)
			os.Exit(1)
		}

		if relatedResetReembed {
			if err := resetAndReembedAllNotes(); err != nil {
				fmt.Fprintf(os.Stderr, "Reset/re-embed failed: %v\n", err)
				os.Exit(1)
			}
		}

		var note models.Note
		if err := globalDB.Conn.First(&note, uint(id)).Error; err != nil {
			fmt.Fprintf(os.Stderr, "Note %d not found: %v\n", id, err)
			os.Exit(1)
		}

		fmt.Printf("=== Note #%d ===\n%s\n", note.ID, note.Content)

		if strings.TrimSpace(note.Content) == "" {
			fmt.Println("(empty content — cannot search for related notes)")
			return
		}

		embedding, err := globalDB.GetNoteEmbedding(note.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch cached embedding: %v\n", err)
			os.Exit(1)
		}
		if embedding == nil {
			if globalEmbedClient == nil {
				fmt.Fprintln(os.Stderr, "No cached embedding and no embedding client configured")
				os.Exit(1)
			}
			embedding, err = globalEmbedClient.Embed(note.Content)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Embedding failed: %v\n", err)
				os.Exit(1)
			}
		}

		results, err := globalDB.SearchRelatedNotes(embedding, 6)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Search failed: %v\n", err)
			os.Exit(1)
		}

		printed := 0
		for _, r := range results {
			if r.ID == note.ID {
				continue
			}
			printed++
			fmt.Printf("\n--- Related #%d (distance %.4f) ---\n%s\n", r.ID, 1, r.Content)
			if printed == 5 {
				break
			}
		}

		if printed == 0 {
			fmt.Println("\n(no related notes found)")
		}
	},
}

func init() {
	RelatedNotesCmd.Flags().BoolVar(
		&relatedResetReembed,
		"reset-reembed",
		false,
		"drop and recreate note_vecs, then re-embed all notes before searching",
	)
}
