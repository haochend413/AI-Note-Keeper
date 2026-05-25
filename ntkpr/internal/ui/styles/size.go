package styles

import "github.com/haochend413/bubbles/v2/table"

// Layout carries all derived dimensions used by the UI.
type Layout struct {
	WindowWidth           int
	WindowHeight          int
	ViewMainContentHeight int
	MainContentHeight     int

	TableWidth       int
	EditorWidth      int
	RecentTableWidth int
	DiffWidth        int

	ThreadBaseHeight int
	BranchBaseHeight int
	NotesBaseHeight  int

	RecentDiffContentHeight int
	RecentDiffBoxHeight     int
	TextAreaHeight          int
	ChangeTableHeight       int

	ThreadColumns []table.Column
	BranchColumns []table.Column
	NoteColumns   []table.Column
	RecentColumns []table.Column
	ChangeColumns []table.Column
}

const (
	borderOverhead          = 8
	tableWidthRatio         = 0.35
	idWidthRatio            = 0.08
	timeWidthRatio          = 0.22
	flagWidthRatio          = 0.15
	countWidthRatio         = 0.07
	nameWidthRatio          = 0.51
	contentWidthRatio       = 0.58
	recentThreadWidthRatio  = 0.13
	recentBranchWidthRatio  = 0.10
	recentNoteWidthRatio    = 0.25
	recentFlagsWidthRatio   = 0.06
	changeTypeWidthRatio    = 0.10
	changeIDWidthRatio      = 0.10
	changeTimeWidthRatio    = 0.25
	changeDescWidthRatio    = 0.40
	statusReservedHeight    = 5
	viewReservedHeight      = 3
	tableMathOffset         = 3
	threadBaseOffset        = -1
	branchBaseOffset        = -1
	notesBaseMultiplier     = 8
	notesBaseOffset         = 1
	notesFocusedExtraHeight = 6
	notesNormalExtraHeight  = 2
	threadFocusedExtra      = 6
	threadSemiExtra         = 2
	branchFocusedExtra      = 6
	boxBorderOverhead       = 2
)

func ComputeLayout(windowWidth, windowHeight int) Layout {
	if windowWidth < 1 {
		windowWidth = 1
	}
	if windowHeight < 1 {
		windowHeight = 1
	}

	availableWidth := windowWidth - borderOverhead
	if availableWidth < 1 {
		availableWidth = 1
	}

	tableWidth := int(float64(availableWidth) * tableWidthRatio)
	if tableWidth < 1 {
		tableWidth = 1
	}

	editorWidth := availableWidth - tableWidth
	if editorWidth < 1 {
		editorWidth = 1
	}

	idWidth := max(4, int(float64(tableWidth)*idWidthRatio))
	timeWidth := max(8, int(float64(tableWidth)*timeWidthRatio))
	flagWidth := max(4, int(float64(tableWidth)*flagWidthRatio))
	countWidth := max(4, int(float64(tableWidth)*countWidthRatio))
	nameWidth := max(10, int(float64(tableWidth)*nameWidthRatio))
	contentWidth := max(10, int(float64(tableWidth)*contentWidthRatio))

	recentThreadWidth := max(20, int(float64(windowWidth)*recentThreadWidthRatio))
	recentBranchWidth := max(20, int(float64(windowWidth)*recentBranchWidthRatio))
	recentNoteWidth := max(20, int(float64(windowWidth)*recentNoteWidthRatio))
	recentFlagsWidth := max(8, int(float64(windowWidth)*recentFlagsWidthRatio))
	recentTableWidth := recentThreadWidth + recentBranchWidth + recentNoteWidth + recentFlagsWidth

	mainContentHeight := windowHeight - statusReservedHeight
	viewMainContentHeight := windowHeight - viewReservedHeight
	tableHeight := max(3, (mainContentHeight-tableMathOffset)/10)
	threadBaseHeight := tableHeight + threadBaseOffset
	branchBaseHeight := tableHeight + branchBaseOffset
	notesBaseHeight := tableHeight*notesBaseMultiplier + notesBaseOffset

	textareaHeight := max(5, mainContentHeight) - 1
	changeTableHeight := max(5, int(float64(mainContentHeight)*0.3)) - 3

	threadColumns := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Name", Width: nameWidth},
		{Title: "#Bs", Width: countWidth},
		{Title: "Flags", Width: flagWidth},
	}

	branchColumns := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Name", Width: nameWidth},
		{Title: "#Ns", Width: countWidth},
		{Title: "Flags", Width: flagWidth},
	}

	noteColumns := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Time", Width: timeWidth},
		{Title: "Content", Width: contentWidth},
		{Title: "Flags", Width: flagWidth},
	}

	recentColumns := []table.Column{
		{Title: "Thread", Width: recentThreadWidth},
		{Title: "Branch", Width: recentBranchWidth},
		{Title: "Note", Width: recentNoteWidth},
		{Title: "Flags", Width: recentFlagsWidth},
	}

	changeColumns := []table.Column{
		{Title: "Type", Width: max(6, int(float64(editorWidth)*changeTypeWidthRatio))},
		{Title: "ID", Width: max(4, int(float64(editorWidth)*changeIDWidthRatio))},
		{Title: "Time", Width: max(12, int(float64(editorWidth)*changeTimeWidthRatio))},
		{Title: "Description", Width: max(15, int(float64(editorWidth)*changeDescWidthRatio))},
	}

	return Layout{
		WindowWidth:             windowWidth,
		WindowHeight:            windowHeight,
		ViewMainContentHeight:   viewMainContentHeight,
		MainContentHeight:       mainContentHeight,
		TableWidth:              tableWidth,
		EditorWidth:             editorWidth,
		RecentTableWidth:        recentTableWidth,
		DiffWidth:               recentTableWidth / 2,
		ThreadBaseHeight:        threadBaseHeight,
		BranchBaseHeight:        branchBaseHeight,
		NotesBaseHeight:         notesBaseHeight,
		RecentDiffContentHeight: notesBaseHeight,
		RecentDiffBoxHeight:     notesBaseHeight + boxBorderOverhead,
		TextAreaHeight:          textareaHeight,
		ChangeTableHeight:       changeTableHeight,
		ThreadColumns:           threadColumns,
		BranchColumns:           branchColumns,
		NoteColumns:             noteColumns,
		RecentColumns:           recentColumns,
		ChangeColumns:           changeColumns,
	}
}

func ThreadFocusedHeight(layout Layout) int {
	return layout.ThreadBaseHeight + threadFocusedExtra
}

func ThreadSemiFocusedHeight(layout Layout) int {
	return layout.ThreadBaseHeight + threadSemiExtra
}

func BranchFocusedHeight(layout Layout) int {
	return layout.BranchBaseHeight + branchFocusedExtra
}

func NotesFocusedHeight(layout Layout) int {
	return layout.NotesBaseHeight + notesFocusedExtraHeight
}

func NotesNormalHeight(layout Layout) int {
	return layout.NotesBaseHeight + notesNormalExtraHeight
}
