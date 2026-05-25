package db

// This package might need further tuning.
// For starters, I think, maybe threads, branches and notes are redundant ?
import (
	"fmt"
	"strings"

	editstack "github.com/haochend413/ntkpr/internal/app/editStack"
	"github.com/haochend413/ntkpr/internal/models"
	"gorm.io/gorm"
)

// SyncData takes in local stored data and edit record, sync with database and return the latest synced data.
func (d *DB) SyncData(
	threads []*models.Thread,
	editMap map[editstack.EditKey]*editstack.Edit) ([]*models.Thread, error) {
	// Categorize edits from the editMap
	noteCreateIDs := make([]uint, 0)
	notePendingIDs := make([]uint, 0)
	noteDeleteIDs := make([]uint, 0)

	threadCreateIDs := make([]uint, 0)
	threadPendingIDs := make([]uint, 0)
	threadDeleteIDs := make([]uint, 0)

	branchCreateIDs := make([]uint, 0)
	branchPendingIDs := make([]uint, 0)
	branchDeleteIDs := make([]uint, 0)

	for key, edit := range editMap {
		id := key.ID
		switch edit.EditType {
		// Notes
		case editstack.CreateNote:
			noteCreateIDs = append(noteCreateIDs, id)
		case editstack.UpdateNote:
			notePendingIDs = append(notePendingIDs, id)
		case editstack.DeleteNote:
			noteDeleteIDs = append(noteDeleteIDs, id)

		// Threads
		case editstack.CreateThread:
			threadCreateIDs = append(threadCreateIDs, id)
		case editstack.UpdateThread:
			threadPendingIDs = append(threadPendingIDs, id)
		case editstack.DeleteThread:
			threadDeleteIDs = append(threadDeleteIDs, id)

		// Branches
		case editstack.CreateBranch:
			branchCreateIDs = append(branchCreateIDs, id)
		case editstack.UpdateBranch:
			branchPendingIDs = append(branchPendingIDs, id)
		case editstack.DeleteBranch:
			branchDeleteIDs = append(branchDeleteIDs, id)

		case editstack.None:
			// Skip
		}
	}

	// I am not sure.
	noteCreateIDs = uniqueIDs(noteCreateIDs)
	notePendingIDs = uniqueIDs(notePendingIDs)
	noteDeleteIDs = uniqueIDs(noteDeleteIDs)
	threadCreateIDs = uniqueIDs(threadCreateIDs)
	threadPendingIDs = uniqueIDs(threadPendingIDs)
	threadDeleteIDs = uniqueIDs(threadDeleteIDs)
	branchCreateIDs = uniqueIDs(branchCreateIDs)
	branchPendingIDs = uniqueIDs(branchPendingIDs)
	branchDeleteIDs = uniqueIDs(branchDeleteIDs)

	// Build maps for O(1) lookup
	threadsMap := make(map[uint]*models.Thread)
	branchesMap := make(map[uint]*models.Branch)
	notesMap := make(map[uint]*models.Note)

	for _, thread := range threads {
		threadsMap[thread.ID] = thread
		for _, branch := range thread.Branches {
			branchesMap[branch.ID] = branch
			for _, note := range branch.Notes {
				notesMap[note.ID] = note
			}
		}
	}

	// Create in order: Threads -> Branches -> Notes
	// Note: nextCreateID logic ensures assigned IDs never collide with existing DB records.
	// SQLite preserves explicitly provided IDs, so foreign key references remain valid.

	// 1. Create threads
	for _, threadID := range threadCreateIDs {
		if thread, exists := threadsMap[threadID]; exists {
			if err := d.persistThread(thread, true); err != nil {
				return nil, fmt.Errorf("failed to create thread %d: %w", thread.ID, err)
			}
		}
	}

	// 2. Create branches
	for _, branchID := range branchCreateIDs {
		if branch, exists := branchesMap[branchID]; exists {
			if err := d.persistBranch(branch, true); err != nil {
				return nil, fmt.Errorf("failed to create branch %d: %w", branch.ID, err)
			}
		}
	}

	// 2.5. Update threads
	for _, threadID := range threadPendingIDs {
		if thread, exists := threadsMap[threadID]; exists {
			if err := d.persistThread(thread, false); err != nil {
				return nil, fmt.Errorf("failed to update thread %d: %w", thread.ID, err)
			}
		}
	}

	// 3. Create notes
	for _, noteID := range noteCreateIDs {
		if note, exists := notesMap[noteID]; exists {
			sanitizeNote(note)
			if err := d.persistNote(note, true); err != nil {
				return nil, fmt.Errorf("failed to create note %d: %w", note.ID, err)
			}
		}
	}

	// 4. Update notes
	for _, noteID := range notePendingIDs {
		if note, exists := notesMap[noteID]; exists {
			sanitizeNote(note)
			if err := d.persistNote(note, false); err != nil {
				return nil, fmt.Errorf("failed to update note %d: %w", note.ID, err)
			}
		}
	}

	// 5. Update branches (e.g., adding/removing notes)
	for _, branchID := range branchPendingIDs {
		if branch, exists := branchesMap[branchID]; exists {
			if err := d.persistBranch(branch, false); err != nil {
				return nil, fmt.Errorf("failed to update branch %d: %w", branch.ID, err)
			}
		}
	}

	// 6. Delete in reverse order: Notes -> Branches -> Threads
	if err := d.deleteNotes(noteDeleteIDs); err != nil {
		return nil, err
	}

	if err := d.deleteBranches(branchDeleteIDs); err != nil {
		return nil, err
	}

	if err := d.deleteThreads(threadDeleteIDs); err != nil {
		return nil, err
	}

	// 7. Load fresh data from DB (with full preloading of hierarchy)
	return d.loadAll()
}

func (d *DB) persistNote(note *models.Note, isCreate bool) error {
	if note == nil {
		return nil
	}
	// Topics removed: only persist note and its branch associations
	var result *gorm.DB
	if isCreate {
		result = d.Conn.Omit("Branches").Create(note)
	} else {
		result = d.Conn.Save(note)
	}
	if result.Error != nil {
		return result.Error
	}

	// Handle branch associations
	return d.Conn.Model(note).Association("Branches").Replace(note.Branches)
}

func (d *DB) deleteNotes(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Note{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func sanitizeNote(note *models.Note) {
	if note == nil {
		return
	}
	note.Content = strings.TrimSpace(note.Content)
}

// func collectNotes(notesMap map[uint]*models.Note, ids []uint) []*models.Note {
// 	if len(ids) == 0 || notesMap == nil {
// 		return nil
// 	}
// 	notes := make([]*models.Note, 0, len(ids))
// 	for _, id := range ids {
// 		if note, exists := notesMap[id]; exists && note != nil {
// 			notes = append(notes, note)
// 		}
// 	}
// 	return notes
// }

func (d *DB) persistThread(thread *models.Thread, isCreate bool) error {
	if thread == nil {
		return nil
	}

	var result *gorm.DB
	if isCreate {
		// When creating a thread, OMIT branches to prevent auto-insert
		// Branches are created separately via their own CreateBranch edits
		result = d.Conn.Omit("Branches").Create(thread)
	} else {
		result = d.Conn.Save(thread)
		if result.Error != nil {
			return result.Error
		}
		// Only replace branch associations when updating
		return d.Conn.Model(thread).Association("Branches").Replace(thread.Branches)
	}

	return result.Error
}

// func collectThreads(threadsMap map[uint]*models.Thread, ids []uint) []*models.Thread {
// 	if len(ids) == 0 || threadsMap == nil {
// 		return nil
// 	}
// 	threads := make([]*models.Thread, 0, len(ids))
// 	for _, id := range ids {
// 		if t, exists := threadsMap[id]; exists && t != nil {
// 			threads = append(threads, t)
// 		}
// 	}
// 	return threads
// }

func (d *DB) deleteThreads(ids []uint) error {
	for _, id := range ids {
		// Cascading delete will handle branches due to foreign key
		// This will automatically delete the related branches.
		if err := d.Conn.Delete(&models.Thread{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) persistBranch(branch *models.Branch, isCreate bool) error {
	if branch == nil {
		return nil
	}

	var result *gorm.DB
	if isCreate {
		// When creating a branch, OMIT notes to prevent auto-insert
		// Notes are created separately via their own CreateNote edits
		result = d.Conn.Omit("Notes").Create(branch)
	} else {
		result = d.Conn.Save(branch)
		if result.Error != nil {
			return result.Error
		}
		// Only replace note associations when updating
		return d.Conn.Model(branch).Association("Notes").Replace(branch.Notes)
	}

	return result.Error
}

func (d *DB) deleteBranches(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Branch{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

// func collectBranches(branchesMap map[uint]*models.Branch, ids []uint) []*models.Branch {
// 	if len(ids) == 0 || branchesMap == nil {
// 		return nil
// 	}
// 	branches := make([]*models.Branch, 0, len(ids))
// 	for _, id := range ids {
// 		if branch, exists := branchesMap[id]; exists && branch != nil {
// 			branches = append(branches, branch)
// 		}
// 	}
// 	return branches
// }

func (d *DB) loadAll() ([]*models.Thread, error) {
	var dbThreads []*models.Thread
	// preload with threads and branches/notes
	if err := d.Conn.
		Preload("Branches.Notes.Branches").
		Order("created_at ASC").
		Find(&dbThreads).Error; err != nil {
		return nil, err
	}

	return dbThreads, nil
}

func uniqueIDs(ids []uint) []uint {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
