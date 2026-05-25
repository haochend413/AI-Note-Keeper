package app

import (
	"log"
	"strings"
	"time"

	editstack "github.com/haochend413/ntkpr/internal/app/editStack"
	"github.com/haochend413/ntkpr/internal/models"
)

// current_branch.go provides a controlled interface for accessing and modifying
// the currently selected branch, with proper edit tracking and synchronization.

// =============================================================================
// Helper: Get current branch with safety checks
// =============================================================================

func (a *App) getCurrentBranch() *models.Branch {
	if a.dataMgr == nil {
		log.Fatal("Critical error: dataMgr is nil - app not properly initialized")
	}
	return a.dataMgr.GetActiveBranch()
}

// =============================================================================
// Getters - Read-only access to current branch properties
// =============================================================================

// GetCurrentBranchID returns the ID of the current branch, or 0 if none selected
func (a *App) GetCurrentBranchID() uint {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return 0
	}
	return branch.ID
}

// GetCurrentBranchName returns the name of the current branch, or empty string if none selected
func (a *App) GetCurrentBranchName() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return ""
	}
	return branch.Name
}

func (a *App) GetCurrentBranchSummary() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return ""
	}
	return branch.Summary
}

// GetCurrentBranchHighlight returns whether the current branch is highlighted
func (a *App) GetCurrentBranchHighlight() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return false
	}
	return branch.Highlight
}

// GetCurrentBranchPrivate returns whether the current branch is private
func (a *App) GetCurrentBranchPrivate() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return false
	}
	return branch.Private
}

// GetCurrentBranchNoteCount returns the number of notes in the current branch
func (a *App) GetCurrentBranchNoteCount() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return 0
	}
	return len(branch.Notes)
}

// GetCurrentBranchNotes returns a copy of the current branch's notes to prevent external mutation.
// This function should not exist. All switching should be handled by the dataMgr.
func (a *App) GetCurrentBranchNotes() []*models.Note {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return nil
	}

	// Return a copy to prevent external mutation
	notes := make([]*models.Note, len(branch.Notes))
	copy(notes, branch.Notes)
	return notes
}

// GetCurrentBranchLastEdit returns the last edit timestamp of the current branch.
func (a *App) GetCurrentBranchLastEdit() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return time.Time{}
	}
	return branch.LastEdit
}

// GetCurrentBranchUpdatedAt returns when the current branch was last modified
func (a *App) GetCurrentBranchUpdatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return time.Time{}
	}
	return branch.UpdatedAt
}

// GetCurrentBranchCreatedAt returns when the current branch was created
func (a *App) GetCurrentBranchCreatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return time.Time{}
	}
	return branch.CreatedAt
}

// IncrementCurrentBranchFrequency increments the current branch's frequency by 1
// and marks it updated for syncing and hooks.
func (a *App) IncrementCurrentBranchFrequency(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	branch.Frequency += 1
	branch.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: branch.ID, EditType: editstack.UpdateBranch}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking branch frequency increment: %v", err)
	}
}

// GetCurrentBranchFrequency returns the edit count (frequency) of the current branch
func (a *App) GetCurrentBranchFrequency() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return 0
	}
	return branch.Frequency
}

// HasCurrentBranch checks if a branch is currently selected
func (a *App) HasCurrentBranch() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.getCurrentBranch() != nil
}

// =============================================================================
// Setters - Controlled modification with edit tracking
// =============================================================================

// SetCurrentBranchName updates the current branch's name with edit tracking
func (a *App) SetCurrentBranchName(name string, link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	// No-op if name hasn't changed
	if branch.Name == name {
		return
	}

	branch.Name = name
	branch.LastEdit = time.Now()
	branch.Frequency += 1
	a.Synced = false

	edit := &editstack.Edit{ID: branch.ID, EditType: editstack.UpdateBranch}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking branch update: %v", err)
	}
}

// SetCurrentBranchSummary updates the current branch's summary with edit tracking
func (a *App) SetCurrentBranchSummary(summary string, link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	// No-op if summary hasn't changed
	if branch.Summary == summary {
		return
	}
	branch.Summary = summary
	lines := strings.Split(summary, "\n")
	if len(lines) > 0 {
		branch.Name = lines[0]
	} else {
		branch.Name = summary
	}
	branch.LastEdit = time.Now()
	branch.Frequency += 1
	a.Synced = false

	edit := &editstack.Edit{ID: branch.ID, EditType: editstack.UpdateBranch}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking branch update: %v", err)
	}
}

// SetCurrentBranchLastEdit updates the LastEdit timestamp of the current branch to the current time.
// Ensures the timestamp is not set to a past time.
func (a *App) SetCurrentBranchLastEdit() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	// dont set it backwards
	if time.Now().Before(branch.LastEdit) {
		return
	}
	branch.LastEdit = time.Now()
}

// ToggleCurrentBranchHighlight toggles the highlight status of the current branch
func (a *App) ToggleCurrentBranchHighlight(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	branch.Highlight = !branch.Highlight
	branch.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: branch.ID, EditType: editstack.UpdateBranch}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking branch update: %v", err)
	}
}

// ToggleCurrentBranchPrivate toggles the private status of the current branch
func (a *App) ToggleCurrentBranchPrivate(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	branch.Private = !branch.Private
	branch.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: branch.ID, EditType: editstack.UpdateBranch}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking branch update: %v", err)
	}
}

// =============================================================================
// Branch Deletion
// =============================================================================

// DeleteCurrentBranch removes the current branch from the current thread and tracks the deletion
func (a *App) DeleteCurrentBranch(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	branch := a.getCurrentBranch()
	if branch == nil {
		return
	}

	branchID := branch.ID

	// Clean up pending CreateNote edits for notes inside this branch
	for _, note := range branch.Notes {
		if noteEdit, noteExists := a.editMgr.GetEdit(editstack.EntityNote, note.ID); noteExists && noteEdit.EditType == editstack.CreateNote {
			a.editMgr.RemoveEdit(editstack.EntityNote, note.ID)
		}
	}

	// Determine if branch was just created or exists in DB
	edit, exists := a.editMgr.GetEdit(editstack.EntityBranch, branchID)

	if exists && edit.EditType == editstack.CreateBranch {
		// Branch was created but not yet synced - just discard it
		a.editMgr.RemoveEdit(editstack.EntityBranch, branchID)
	} else if branchID != 0 {
		// Branch exists in DB - mark for deletion
		deleteEdit := &editstack.Edit{ID: branchID, EditType: editstack.DeleteBranch}
		if err := a.editMgr.AddEdit(deleteEdit, link); err != nil {
			log.Printf("Error tracking branch deletion: %v", err)
			return
		}
	}

	// Find the index of the branch in the active branch list
	branches := a.dataMgr.GetActiveBranchList()
	for i, b := range branches {
		if b.ID == branchID {
			a.dataMgr.RemoveBranch(i)
			break
		}
	}

	a.Synced = false
}
