package ui

import (
	"github.com/haochend413/ntkpr/internal/app/context"
	"github.com/haochend413/ntkpr/state"
)

// DistributeState restores cursor positions from saved state on startup.
func (m *Model) DistributeState(s *state.AppState) {
	if s == nil {
		return
	}
	if tc, ok := s.ThreadCursors[context.Default]; ok {
		if int(tc) < len(m.threadsTable.Rows()) {
			m.threadsTable.SetCursor(int(tc))
			m.switchToThreadAtCursor(int(tc))
			m.updateBranchesTable()
		}
	}
	if bc, ok := s.BranchCursors[context.Default]; ok {
		if int(bc) < len(m.branchesTable.Rows()) {
			m.branchesTable.SetCursor(int(bc))
			m.switchToBranchAtCursor(int(bc))
			m.updateNotesTable()
		}
	}
	if nc, ok := s.NoteCursors[context.Default]; ok {
		if int(nc) < len(m.notesTable.Rows()) {
			m.notesTable.SetCursor(int(nc))
			m.switchToNoteAtCursor(int(nc))
		}
	}
}

// CollectState gathers current cursor positions for persistence on quit.
func (m Model) CollectState() *state.State {
	s := state.DefaultState()
	s.App.ThreadCursors[context.Default] = uint(m.threadsTable.Cursor())
	s.App.BranchCursors[context.Default] = uint(m.branchesTable.Cursor())
	s.App.NoteCursors[context.Default] = uint(m.notesTable.Cursor())
	return s
}

func HelpText() string {
	help :=
		`ntkpr is a note / diary management tool built with Golang and Bubbletea framework. 

It allows you to manage your logs in a structured manner. The general idea comes from version control tools in software development. 
You can separate your tasks into different threads, and within one thread you can create multiple branches, each containing a list of notes. 
Each branch / thread contains one summary page that you can use to provide info for this specific branch / thread. 

I recommand creating several short notes instead of one long note. 

**********************

Some important keybindings: 

q : switch to upper table 
w : switch to lower table 
tab : switch between windows 
e : go to edit window 
n : create new item 
ctrl-s : save current edit 

For full usage guidance, please check github : https://github.com/haochend413/ntkpr	
	`
	return help
}
