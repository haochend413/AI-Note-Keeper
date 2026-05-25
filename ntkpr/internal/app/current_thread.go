package app

import (
	"log"
	"strings"
	"time"

	editstack "github.com/haochend413/ntkpr/internal/app/editStack"
	"github.com/haochend413/ntkpr/internal/models"
)

// current_thread.go provides a controlled interface for accessing and modifying
// the currently selected thread, with proper edit tracking and synchronization.

// =============================================================================
// Helper: Get current thread with safety checks
// =============================================================================

func (a *App) getCurrentThread() *models.Thread {
	if a.dataMgr == nil {
		log.Fatal("Critical error: dataMgr is nil - app not properly initialized")
	}
	return a.dataMgr.GetActiveThread()
}

// =============================================================================
// Getters - Read-only access to current thread properties
// =============================================================================

// GetCurrentThreadID returns the ID of the current thread, or 0 if none selected
func (a *App) GetCurrentThreadID() uint {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return 0
	}
	return thread.ID
}

// GetCurrentThreadName returns the name of the current thread, or empty string if none selected
func (a *App) GetCurrentThreadName() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return ""
	}
	return thread.Name
}

func (a *App) GetCurrentThreadSummary() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return ""
	}
	return thread.Summary
}

// GetCurrentThreadHighlight returns whether the current thread is highlighted
func (a *App) GetCurrentThreadHighlight() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return false
	}
	return thread.Highlight
}

// GetCurrentThreadPrivate returns whether the current thread is private
func (a *App) GetCurrentThreadPrivate() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return false
	}
	return thread.Private
}

// GetCurrentThreadBranchCount returns the number of branches in the current thread
func (a *App) GetCurrentThreadBranchCount() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return 0
	}
	return len(thread.Branches)
}

// GetCurrentThreadBranches returns a copy of the current thread's branches to prevent external mutation
// This function should not exist. All switching should be handled by the dataMgr.
func (a *App) GetCurrentThreadBranches() []*models.Branch {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return nil
	}

	// Return a copy to prevent external mutation
	branches := make([]*models.Branch, len(thread.Branches))
	copy(branches, thread.Branches)
	return branches
}

// GetCurrentThreadLastEdit returns the last edit timestamp of the current thread.
func (a *App) GetCurrentThreadLastEdit() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return time.Time{}
	}
	return thread.LastEdit
}

// GetCurrentThreadUpdatedAt returns when the current thread was last modified
func (a *App) GetCurrentThreadUpdatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return time.Time{}
	}
	return thread.UpdatedAt
}

// GetCurrentThreadCreatedAt returns when the current thread was created
func (a *App) GetCurrentThreadCreatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return time.Time{}
	}
	return thread.CreatedAt
}

// IncrementCurrentThreadFrequency increments the current thread's frequency by 1
// and marks it updated for syncing and hooks.
func (a *App) IncrementCurrentThreadFrequency(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	thread.Frequency += 1
	thread.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: thread.ID, EditType: editstack.UpdateThread}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking thread frequency increment: %v", err)
	}
}

// GetCurrentThreadFrequency returns the edit count (frequency) of the current thread
func (a *App) GetCurrentThreadFrequency() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return 0
	}
	return thread.Frequency
}

// HasCurrentThread checks if a thread is currently selected
func (a *App) HasCurrentThread() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.getCurrentThread() != nil
}

// =============================================================================
// Setters - Controlled modification with edit tracking
// =============================================================================

// SetCurrentThreadName updates the current thread's name with edit tracking
func (a *App) SetCurrentThreadName(name string, link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	// No-op if name hasn't changed
	if thread.Name == name {
		return
	}

	thread.Name = name
	thread.LastEdit = time.Now()
	thread.Frequency += 1
	a.Synced = false

	edit := &editstack.Edit{ID: thread.ID, EditType: editstack.UpdateThread}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking thread update: %v", err)
	}
}

// SetCurrentThreadSummary updates the current thread's summary with edit tracking
func (a *App) SetCurrentThreadSummary(summary string, link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	// No-op if summary hasn't changed
	if thread.Summary == summary {
		return
	}

	thread.Summary = summary
	lines := strings.Split(summary, "\n")
	if len(lines) > 0 {
		thread.Name = lines[0]
	} else {
		thread.Name = summary
	}
	thread.LastEdit = time.Now()
	thread.Frequency += 1
	a.Synced = false

	edit := &editstack.Edit{ID: thread.ID, EditType: editstack.UpdateThread}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking thread update: %v", err)
	}
}

// SetCurrentThreadLastEdit updates the LastEdit timestamp of the current thread to the current time.
// Ensures the timestamp is not set to a past time.
func (a *App) SetCurrentThreadLastEdit() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	// dont set it backwards
	if time.Now().Before(thread.LastEdit) {
		return
	}
	thread.LastEdit = time.Now()
}

// ToggleCurrentThreadHighlight toggles the highlight status of the current thread
func (a *App) ToggleCurrentThreadHighlight(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	thread.Highlight = !thread.Highlight
	thread.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: thread.ID, EditType: editstack.UpdateThread}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking thread update: %v", err)
	}
}

// ToggleCurrentThreadPrivate toggles the private status of the current thread
func (a *App) ToggleCurrentThreadPrivate(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	thread.Private = !thread.Private
	thread.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: thread.ID, EditType: editstack.UpdateThread}
	if err := a.editMgr.AddEdit(edit, link); err != nil {
		log.Printf("Error tracking thread update: %v", err)
	}
}

// =============================================================================
// Thread Deletion
// =============================================================================

// DeleteCurrentThread removes the current thread and all its branches and tracks the deletion
func (a *App) DeleteCurrentThread(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.getCurrentThread()
	if thread == nil {
		return
	}

	threadID := thread.ID

	// Clean up pending CreateBranch/CreateNote edits for entities inside this thread
	for _, branch := range thread.Branches {
		if branchEdit, branchExists := a.editMgr.GetEdit(editstack.EntityBranch, branch.ID); branchExists && branchEdit.EditType == editstack.CreateBranch {
			a.editMgr.RemoveEdit(editstack.EntityBranch, branch.ID)
		}
		for _, note := range branch.Notes {
			if noteEdit, noteExists := a.editMgr.GetEdit(editstack.EntityNote, note.ID); noteExists && noteEdit.EditType == editstack.CreateNote {
				a.editMgr.RemoveEdit(editstack.EntityNote, note.ID)
			}
		}
	}

	// Determine if thread was just created or exists in DB
	edit, exists := a.editMgr.GetEdit(editstack.EntityThread, threadID)

	if exists && edit.EditType == editstack.CreateThread {
		// Thread was created but not yet synced - just discard it
		a.editMgr.RemoveEdit(editstack.EntityThread, threadID)
	} else if threadID != 0 {
		// Thread exists in DB - mark for deletion (will cascade to branches)
		deleteEdit := &editstack.Edit{ID: threadID, EditType: editstack.DeleteThread}
		if err := a.editMgr.AddEdit(deleteEdit, link); err != nil {
			log.Printf("Error tracking thread deletion: %v", err)
			return
		}
	}

	// Find the index of the thread in the thread list
	threads := a.dataMgr.GetThreads()
	for i, t := range threads {
		if t.ID == threadID {
			a.dataMgr.RemoveThread(i)
			break
		}
	}

	a.Synced = false
}
