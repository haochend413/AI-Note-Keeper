package data

import "github.com/haochend413/ntkpr/internal/models"

// DataMgr should handle the switching logic between threads, branches and notes. It keeps record of all threads, and exposing current threads, branches and notes.
// DataMgr should only expose 1 threadlist, 1 branchlist, and 1 notelist. ContextMgr will Demonstrate based on that. THe optimization and storage should happen at this level.
// Maybe a context manager should have a DataMgr inside? I suppose.
// DataMgr will also be responsible for syncing with database, instead of being handled directly by App.
// It is the one-above layer of db.
// With this package, we can easily manage what local data is stored for our app.
// Cursor needs much more careful handling!
// I came up with this at Houston Airport IAH. What a stroke of genious.
// I will have to re-write the context handling stuff based on this. Awwww
type DataMgr struct {
	threads         []*models.Thread // All the threads
	branches        []*models.Branch // The branch list of the active thread
	notes           []*models.Note   // the note list of the active branch
	activeThreadPtr int              // which thread is active ? // Maybe this should be passed down.
	activeBranchPtr int              //which branch is active ? // Maybe this should be passed down.
	activeNotePtr   int              //This is less useful, just a cursor.
	activeThreadID  uint
	activeBranchID  uint
	activeNoteID    uint
	threadIndexByID map[uint]int
	branchIndexByID map[uint]int
	noteIndexByID   map[uint]int
	// Global maps for O(1) lookup across all threads/branches
	globalThreadByID map[uint]*models.Thread
	globalBranchByID map[uint]*models.Branch
	globalNoteByID   map[uint]*models.Note
}

func NewDataMgr(threads []*models.Thread) *DataMgr {
	dm := &DataMgr{
		threads:          threads,
		activeThreadPtr:  0,
		activeBranchPtr:  0,
		activeNotePtr:    0,
		threadIndexByID:  make(map[uint]int),
		branchIndexByID:  make(map[uint]int),
		noteIndexByID:    make(map[uint]int),
		globalThreadByID: make(map[uint]*models.Thread),
		globalBranchByID: make(map[uint]*models.Branch),
		globalNoteByID:   make(map[uint]*models.Note),
	}

	dm.rebuildThreadIndex()

	// Initialize branches and notes from first thread if available
	if len(threads) > 0 {
		dm.branches = threads[0].Branches
		dm.activeThreadID = threads[0].ID
		dm.rebuildBranchIndex()
		if len(dm.branches) > 0 {
			dm.notes = dm.branches[0].Notes
			dm.activeBranchID = dm.branches[0].ID
			dm.rebuildNoteIndex()
			if len(dm.notes) > 0 {
				dm.activeNoteID = dm.notes[0].ID
			}
		} else {
			dm.notes = []*models.Note{}
			dm.activeBranchID = 0
			dm.activeNoteID = 0
		}
	} else {
		dm.branches = []*models.Branch{}
		dm.notes = []*models.Note{}
		dm.activeThreadID = 0
		dm.activeBranchID = 0
		dm.activeNoteID = 0
	}

	return dm
}

func (dm *DataMgr) NewDataMgr() *DataMgr {
	return &DataMgr{
		threads:          []*models.Thread{},
		branches:         []*models.Branch{},
		notes:            []*models.Note{},
		threadIndexByID:  make(map[uint]int),
		branchIndexByID:  make(map[uint]int),
		noteIndexByID:    make(map[uint]int),
		globalThreadByID: make(map[uint]*models.Thread),
		globalBranchByID: make(map[uint]*models.Branch),
		globalNoteByID:   make(map[uint]*models.Note),
	}
}

func (dm *DataMgr) rebuildThreadIndex() {
	dm.threadIndexByID = make(map[uint]int, len(dm.threads))
	for i, t := range dm.threads {
		dm.threadIndexByID[t.ID] = i
	}
	dm.rebuildGlobalMaps()
}

func (dm *DataMgr) rebuildGlobalMaps() {
	dm.globalThreadByID = make(map[uint]*models.Thread, len(dm.threads))
	totalBranches := 0
	for _, t := range dm.threads {
		totalBranches += len(t.Branches)
	}
	dm.globalBranchByID = make(map[uint]*models.Branch, totalBranches)
	dm.globalNoteByID = make(map[uint]*models.Note)
	for _, t := range dm.threads {
		dm.globalThreadByID[t.ID] = t
		for _, b := range t.Branches {
			dm.globalBranchByID[b.ID] = b
			for _, n := range b.Notes {
				dm.globalNoteByID[n.ID] = n
			}
		}
	}
}

func (dm *DataMgr) rebuildBranchIndex() {
	dm.branchIndexByID = make(map[uint]int, len(dm.branches))
	for i, b := range dm.branches {
		dm.branchIndexByID[b.ID] = i
	}
}

func (dm *DataMgr) rebuildNoteIndex() {
	dm.noteIndexByID = make(map[uint]int, len(dm.notes))
	for i, n := range dm.notes {
		dm.noteIndexByID[n.ID] = i
	}
}

func (dm *DataMgr) resetBranchAndNoteState() {
	dm.branches = []*models.Branch{}
	dm.activeBranchPtr = 0
	dm.activeBranchID = 0
	dm.branchIndexByID = make(map[uint]int)
	dm.notes = []*models.Note{}
	dm.activeNotePtr = 0
	dm.activeNoteID = 0
	dm.noteIndexByID = make(map[uint]int)
}

func (dm *DataMgr) GetThreads() []*models.Thread {
	return dm.threads
}

func (dm *DataMgr) GetActiveThread() *models.Thread {
	if len(dm.threads) == 0 || dm.activeThreadPtr < 0 || dm.activeThreadPtr >= len(dm.threads) {
		return nil
	}
	return dm.threads[dm.activeThreadPtr]
}

func (dm *DataMgr) GetActiveBranchList() []*models.Branch {
	return dm.branches
}

func (dm *DataMgr) GetActiveBranch() *models.Branch {
	if len(dm.branches) == 0 || dm.activeBranchPtr < 0 || dm.activeBranchPtr >= len(dm.branches) {
		return nil
	}
	return dm.branches[dm.activeBranchPtr]
}

func (dm *DataMgr) GetActiveNoteList() []*models.Note {
	return dm.notes
}

func (dm *DataMgr) GetActiveNote() *models.Note {
	if len(dm.notes) == 0 || dm.activeNotePtr < 0 || dm.activeNotePtr >= len(dm.notes) {
		return nil
	}
	return dm.notes[dm.activeNotePtr]
}

// GetActiveThreadPtr returns the current thread pointer
func (dm *DataMgr) GetActiveThreadPtr() int {
	return dm.activeThreadPtr
}

// GetActiveThreadID returns the current thread ID.
func (dm *DataMgr) GetActiveThreadID() uint {
	return dm.activeThreadID
}

// GetActiveBranchPtr returns the current branch pointer
func (dm *DataMgr) GetActiveBranchPtr() int {
	return dm.activeBranchPtr
}

// GetActiveBranchID returns the current branch ID.
func (dm *DataMgr) GetActiveBranchID() uint {
	return dm.activeBranchID
}

// GetActiveNotePtr returns the current note pointer
func (dm *DataMgr) GetActiveNotePtr() int {
	return dm.activeNotePtr
}

// GetActiveNoteID returns the current note ID.
func (dm *DataMgr) GetActiveNoteID() uint {
	return dm.activeNoteID
}

// RefreshData updates datamgr with new thread list.
// This should come with states. Implement later.
func (dm *DataMgr) RefreshData(threads []*models.Thread, tc *int, bc *int, nc *int) {
	dm.threads = threads
	dm.rebuildThreadIndex()
	if tc == nil {
		dm.activeThreadPtr = 0
	} else {
		dm.activeThreadPtr = *tc
	}

	// Handle empty threads or out of bounds
	if len(dm.threads) == 0 || dm.activeThreadPtr < 0 || dm.activeThreadPtr >= len(dm.threads) {
		dm.activeThreadPtr = 0
		dm.activeThreadID = 0
		dm.resetBranchAndNoteState()
		return
	}
	dm.activeThreadID = dm.threads[dm.activeThreadPtr].ID

	dm.branches = dm.threads[dm.activeThreadPtr].Branches
	dm.rebuildBranchIndex()

	if bc == nil {
		dm.activeBranchPtr = 0
	} else {
		dm.activeBranchPtr = *bc
	}

	// Handle empty branches or out of bounds
	if len(dm.branches) == 0 || dm.activeBranchPtr < 0 || dm.activeBranchPtr >= len(dm.branches) {
		dm.activeBranchPtr = 0
		dm.activeBranchID = 0
		dm.notes = []*models.Note{}
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		dm.noteIndexByID = make(map[uint]int)
		return
	}
	dm.activeBranchID = dm.branches[dm.activeBranchPtr].ID

	dm.notes = dm.branches[dm.activeBranchPtr].Notes
	dm.rebuildNoteIndex()

	if nc == nil {
		dm.activeNotePtr = 0
	} else {
		dm.activeNotePtr = *nc
	}

	if len(dm.notes) == 0 || dm.activeNotePtr < 0 || dm.activeNotePtr >= len(dm.notes) {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return
	}

	dm.activeNoteID = dm.notes[dm.activeNotePtr].ID
}

// RefreshDataByID updates datamgr with new thread list and restores active entities by ID.
func (dm *DataMgr) RefreshDataByID(threads []*models.Thread, threadID *uint, branchID *uint, noteID *uint) {
	dm.threads = threads
	dm.rebuildThreadIndex()

	if len(dm.threads) == 0 {
		dm.activeThreadPtr = 0
		dm.activeThreadID = 0
		dm.resetBranchAndNoteState()
		return
	}

	if threadID != nil && dm.SwitchActiveThreadByID(*threadID) {
		if branchID != nil {
			dm.SwitchActiveBranchByID(*branchID)
		}
		if noteID != nil {
			dm.SwitchActiveNoteByID(*noteID)
		}
		return
	}

	dm.SwitchActiveThread(0)
	if branchID != nil {
		dm.SwitchActiveBranchByID(*branchID)
	}
	if noteID != nil {
		dm.SwitchActiveNoteByID(*noteID)
	}
}

// SwitchActiveThread deals with switching threads. It updates the exposed branch list when we switch threads.
func (dm *DataMgr) SwitchActiveThread(cursor int) {
	if len(dm.threads) == 0 {
		dm.activeThreadPtr = 0
		dm.activeThreadID = 0
		dm.resetBranchAndNoteState()
		return
	}

	if cursor < 0 || cursor >= len(dm.threads) {
		return
	}

	_ = dm.SwitchActiveThreadByID(dm.threads[cursor].ID)
}

// SwitchActiveThreadByID switches active thread by stable thread ID.
func (dm *DataMgr) SwitchActiveThreadByID(threadID uint) bool {
	if len(dm.threads) == 0 {
		dm.activeThreadPtr = 0
		dm.activeThreadID = 0
		dm.resetBranchAndNoteState()
		return false
	}

	idx, ok := dm.threadIndexByID[threadID]
	if !ok {
		return false
	}

	dm.activeThreadPtr = idx
	dm.activeThreadID = threadID
	dm.branches = dm.threads[idx].Branches
	dm.rebuildBranchIndex()

	if len(dm.branches) == 0 {
		dm.activeBranchPtr = 0
		dm.activeBranchID = 0
		dm.notes = []*models.Note{}
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		dm.noteIndexByID = make(map[uint]int)
		return true
	}

	dm.activeBranchPtr = 0
	dm.activeBranchID = dm.branches[0].ID
	dm.notes = dm.branches[0].Notes
	dm.rebuildNoteIndex()

	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return true
	}

	dm.activeNotePtr = 0
	dm.activeNoteID = dm.notes[0].ID
	return true
}

// SwitchActiveBranch switches to a different branch within the current thread and resets the note list.
func (dm *DataMgr) SwitchActiveBranch(cursor int) {
	if len(dm.branches) == 0 {
		dm.activeBranchPtr = 0
		dm.activeBranchID = 0
		dm.notes = []*models.Note{}
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		dm.noteIndexByID = make(map[uint]int)
		return
	}

	if cursor < 0 || cursor >= len(dm.branches) {
		return
	}

	_ = dm.SwitchActiveBranchByID(dm.branches[cursor].ID)
}

// SwitchActiveBranchByID switches active branch by stable branch ID.
func (dm *DataMgr) SwitchActiveBranchByID(branchID uint) bool {
	if len(dm.branches) == 0 {
		dm.activeBranchPtr = 0
		dm.activeBranchID = 0
		dm.notes = []*models.Note{}
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		dm.noteIndexByID = make(map[uint]int)
		return false
	}

	idx, ok := dm.branchIndexByID[branchID]
	if !ok {
		return false
	}

	dm.activeBranchPtr = idx
	dm.activeBranchID = branchID
	dm.notes = dm.branches[idx].Notes
	dm.rebuildNoteIndex()

	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return true
	}

	dm.activeNotePtr = 0
	dm.activeNoteID = dm.notes[0].ID
	return true
}

// SwitchActiveNote switches to a different note within the current branch.
func (dm *DataMgr) SwitchActiveNote(cursor int) {
	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return
	}
	if cursor < 0 || cursor >= len(dm.notes) {
		return
	}
	_ = dm.SwitchActiveNoteByID(dm.notes[cursor].ID)
}

// SwitchActiveNoteByID switches active note by stable note ID.
func (dm *DataMgr) SwitchActiveNoteByID(noteID uint) bool {
	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return false
	}

	idx, ok := dm.noteIndexByID[noteID]
	if !ok {
		return false
	}

	dm.activeNotePtr = idx
	dm.activeNoteID = noteID
	return true
}

// AddThread adds a thread to thread list without switching to it.
func (dm *DataMgr) AddThread(t *models.Thread) {
	if t == nil {
		return
	}
	dm.threads = append(dm.threads, t)
	dm.threadIndexByID[t.ID] = len(dm.threads) - 1
	dm.globalThreadByID[t.ID] = t
}

// RemoveThread removes a thread at the given index and adjusts active pointers.
func (dm *DataMgr) RemoveThread(index int) {
	if index < 0 || index >= len(dm.threads) {
		return
	}

	removedThread := dm.threads[index]
	removedID := removedThread.ID
	prevThreadID := dm.activeThreadID

	// Remove from global maps
	delete(dm.globalThreadByID, removedID)
	for _, b := range removedThread.Branches {
		delete(dm.globalBranchByID, b.ID)
		for _, n := range b.Notes {
			delete(dm.globalNoteByID, n.ID)
		}
	}

	dm.threads = append(dm.threads[:index], dm.threads[index+1:]...)
	dm.rebuildThreadIndex()

	if len(dm.threads) == 0 {
		dm.activeThreadPtr = 0
		dm.activeThreadID = 0
		dm.resetBranchAndNoteState()
		return
	}

	if prevThreadID != 0 && prevThreadID != removedID && dm.SwitchActiveThreadByID(prevThreadID) {
		return
	}

	if index >= len(dm.threads) {
		index = len(dm.threads) - 1
	}
	dm.SwitchActiveThread(index)
}

// AddBranch adds a branch to the current thread's branch list without switching to it.
func (dm *DataMgr) AddBranch(b *models.Branch) {
	if b == nil || len(dm.threads) == 0 || dm.activeThreadPtr >= len(dm.threads) {
		return
	}

	thread := dm.threads[dm.activeThreadPtr]
	thread.Branches = append(thread.Branches, b)
	dm.branches = thread.Branches
	dm.branchIndexByID[b.ID] = len(dm.branches) - 1
	dm.globalBranchByID[b.ID] = b
}

// RemoveBranch removes a branch at the given index from the current thread and adjusts active pointers.
func (dm *DataMgr) RemoveBranch(index int) {
	if len(dm.threads) == 0 || index < 0 || index >= len(dm.branches) || dm.activeThreadPtr >= len(dm.threads) {
		return
	}

	removedBranch := dm.branches[index]
	removedID := removedBranch.ID
	prevBranchID := dm.activeBranchID

	// Remove from global maps
	delete(dm.globalBranchByID, removedID)
	for _, n := range removedBranch.Notes {
		delete(dm.globalNoteByID, n.ID)
	}

	thread := dm.threads[dm.activeThreadPtr]
	thread.Branches = append(thread.Branches[:index], thread.Branches[index+1:]...)
	dm.branches = thread.Branches
	dm.rebuildBranchIndex()

	if len(dm.branches) == 0 {
		dm.activeBranchPtr = 0
		dm.activeBranchID = 0
		dm.notes = []*models.Note{}
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		dm.noteIndexByID = make(map[uint]int)
		return
	}

	if prevBranchID != 0 && prevBranchID != removedID && dm.SwitchActiveBranchByID(prevBranchID) {
		return
	}

	if index >= len(dm.branches) {
		index = len(dm.branches) - 1
	}
	dm.SwitchActiveBranch(index)
}

// AddNote adds a note to the current branch's note list without switching to it.
func (dm *DataMgr) AddNote(n *models.Note) {
	if n == nil || len(dm.branches) == 0 || dm.activeBranchPtr >= len(dm.branches) {
		return
	}

	branch := dm.branches[dm.activeBranchPtr]
	branch.Notes = append(branch.Notes, n)
	dm.notes = branch.Notes
	dm.noteIndexByID[n.ID] = len(dm.notes) - 1
	dm.globalNoteByID[n.ID] = n
}

// RemoveNote removes a note at the given index from the current branch and adjusts active pointers.
func (dm *DataMgr) RemoveNote(index int) {
	if len(dm.branches) == 0 || index < 0 || index >= len(dm.notes) {
		return
	}

	removedID := dm.notes[index].ID
	prevNoteID := dm.activeNoteID

	// Remove from global map
	delete(dm.globalNoteByID, removedID)

	branch := dm.branches[dm.activeBranchPtr]
	branch.Notes = append(branch.Notes[:index], branch.Notes[index+1:]...)
	dm.notes = branch.Notes
	dm.rebuildNoteIndex()

	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return
	}

	if prevNoteID != 0 && prevNoteID != removedID && dm.SwitchActiveNoteByID(prevNoteID) {
		return
	}

	if index >= len(dm.notes) {
		index = len(dm.notes) - 1
	}
	dm.SwitchActiveNote(index)
}

// FindThreadByID finds a thread by ID across all threads in O(1).
func (dm *DataMgr) FindThreadByID(id uint) *models.Thread {
	return dm.globalThreadByID[id]
}

// FindBranchByID finds a branch by ID across all threads in O(1).
func (dm *DataMgr) FindBranchByID(id uint) *models.Branch {
	return dm.globalBranchByID[id]
}

// FindNoteByID finds a note by ID across all threads and branches in O(1).
func (dm *DataMgr) FindNoteByID(id uint) *models.Note {
	return dm.globalNoteByID[id]
}

// FindNoteByLink finds a note by superlink in O(1).
func (dm *DataMgr) FindNoteByLink(link models.Superlink) *models.Note {
	if link.NoteID <= 0 {
		return nil
	}
	return dm.globalNoteByID[uint(link.NoteID)]
}
