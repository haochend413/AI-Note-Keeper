package ui

import (
	"fmt"
	"strconv"
	"time"

	// "github.com/haochend413/bubbles/table"
	// "github.com/haochend413/bubbles/textarea_vim"
	// "github.com/haochend413/bubbles/textinput"
	// "github.com/haochend413/bubbles/v2/statusbar"

	// "charm.land/bubbles/table"
	// "charm.land/bubbles/textinput"
	"github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/bubbles/v2/viewport"

	// "github.com/haochend413/bubbles/textarea_vim"
	"github.com/haochend413/bubbles/v2/textinput"

	// "charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	// "charm.land/lipgloss/v2"

	// "github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/v2/statusbar"
	// "charm.land/bubbles/v2/table"
	// "charm.land/bubbles/v2/textinput"

	"github.com/haochend413/bubbles/v2/textarea_vim"
	// "charm.land/bubbles/v2/textarea_vim"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/ntkpr/config"
	"github.com/haochend413/ntkpr/internal/app"
	"github.com/haochend413/ntkpr/internal/models"
	"github.com/haochend413/ntkpr/state"
	"github.com/haochend413/ntkpr/sys"
)

// FocusState represents the current UI focus
type FocusState int

const (
	FocusThreads FocusState = iota
	FocusBranches
	FocusNotes
	FocusEdit
	FocusChangelog
	FocusRecent
	FocusDiff
)

type ViewMode int

const (
	ApplicationView ViewMode = iota
	QuitConfirmView
)

// tickMsg is used to update the UI clock every second.
type tickMsg time.Time

// Model represents the Bubble Tea model
// This should be managed better. At least, something should be rendered into groups.
type Model struct {
	// metadata
	app    *app.App
	Config *config.Config

	// windows
	threadsTable  table.Model
	branchesTable table.Model
	notesTable    table.Model
	textArea      textarea_vim.Model
	changeTable   table.Model
	recentTable   table.Model
	diffView      viewport.Model // we might need something better for this.
	statusBar     statusbar.Model

	//view mode
	viewMode ViewMode

	//states
	previousFocus   FocusState
	focus           FocusState
	editPrevIMEType sys.InputMethodType
	ready           bool

	//data
	width  int
	height int
}

func colorPtr(c string) *string {
	return &c
}

// NewModel initializes a new UI model
func NewModel(application *app.App, cfg *config.Config, s *state.State) Model {
	// Use default state if nil
	// ignore this for now
	if s == nil {
		s = state.DefaultState()
	}
	if cfg == nil {
		temp := config.LoadOrCreateConfig()
		cfg = &temp
	}

	recentColumns := []table.Column{
		{Title: "Thread", Width: 50},
		{Title: "Branch", Width: 50},
		{Title: "Note", Width: 70},
		{Title: "Flags", Width: 16},
	}

	recentTable := table.New(
		table.WithColumns(recentColumns),
		table.WithFocused(true),
		table.WithHeight(40),
	)

	noteColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Content", Width: 40},
		{Title: "Flags", Width: 6},
	}

	noteTable := table.New(
		table.WithColumns(noteColumns),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	branchColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Name", Width: 40},
		{Title: "#Ns", Width: 2},
		{Title: "Flags", Width: 6},
	}

	branchTable := table.New(
		table.WithColumns(branchColumns),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	threadColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Name", Width: 40},
		{Title: "#Bs", Width: 2},
		{Title: "Flags", Width: 6},
	}

	threadTable := table.New(
		table.WithColumns(threadColumns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	//set sb
	//Left: Context ; NoteID ; Last Update ; Version (frequency)
	//Right: Action ; Synced ? ; Time
	// Example usage with method chaining
	sb := statusbar.New(
		statusbar.WithHeight(1),
		statusbar.WithWidth(100),
		statusbar.WithLeftLen(4),
		statusbar.WithRightLen(3),
	)

	// Configure all left elements in sequence
	sb.GetLeft(0).SetValue("Context: Default").SetColors(colorPtr("252"), colorPtr("237")).SetWidth(25)

	sb.GetLeft(1).SetValue("NoteID: -").SetColors(colorPtr("250"), colorPtr("238")).SetWidth(15)
	sb.GetLeft(2).SetValue("Updated: Never").SetColors(colorPtr("250"), colorPtr("239")).SetWidth(20)
	sb.GetLeft(3).SetValue("Version: 1.0").SetColors(colorPtr("250"), colorPtr("240")).SetWidth(15)
	//set tags for quick and consistent access
	sb.SetTag(sb.GetLeft(0), "filter")
	sb.SetTag(sb.GetLeft(1), "ID")
	sb.SetTag(sb.GetLeft(2), "LastUpdated")
	sb.SetTag(sb.GetLeft(3), "Frequency")

	// Configure all right elements in sequence
	sb.GetRight(0).SetValue("").SetColors(colorPtr("226"), nil).SetWidth(30)
	sb.GetRight(1).SetValue("Synced").SetColors(colorPtr("232"), colorPtr("118")).SetWidth(15)
	sb.GetRight(2).SetValue(time.Now().Format("15:04:05")).SetColors(colorPtr("250"), colorPtr("236")).SetWidth(10)
	sb.SetTag(sb.GetRight(0), "Action")
	sb.SetTag(sb.GetRight(1), "Synced")
	sb.SetTag(sb.GetRight(2), "Time")

	// You can also chain model methods
	sb.SetWidth(100).SetHeight(1)

	textArea := textarea_vim.New()

	// Set colors - get styles, modify, and set back
	styles := textArea.Styles()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	styles.Focused.CursorLine = cursorStyle
	styles.Blurred.CursorLine = cursorStyle
	textArea.SetStyles(styles)

	textArea.Placeholder = "Start writing! For summary of thread / branch, the first line of this textarea will be assigned to Name entry." // This should change when we switch between threads / branches / lists
	textArea.SetWidth(50)
	textArea.SetHeight(10)

	//diff
	diffView := viewport.New()

	// diffstyles := diffView.Styles()
	// diffcursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	// diffstyles.Focused.CursorLine = diffcursorStyle
	// diffstyles.Blurred.CursorLine = diffcursorStyle
	// difftextArea.SetStyles(diffstyles)

	// difftextArea.Placeholder = "Start writing! For summary of thread / branch, the first line of this textarea will be assigned to Name entry." // This should change when we switch between threads / branches / lists
	diffView.SetWidth(50)
	diffView.SetHeight(10)
	diffView.SoftWrap = true
	// diffView.GotoBottom()
	// This needs further improving.
	changeColumns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Content", Width: 40},
		{Title: "Flags", Width: 6},
	}

	changeTable := table.New(
		table.WithColumns(changeColumns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	m := Model{
		app:             application,
		Config:          cfg,
		threadsTable:    threadTable,
		branchesTable:   branchTable,
		notesTable:      noteTable,
		recentTable:     recentTable,
		textArea:        textArea,
		diffView:        diffView,
		viewMode:        ApplicationView,
		statusBar:       sb,
		changeTable:     changeTable,
		focus:           FocusThreads,
		editPrevIMEType: sys.InputMethodEnglish, // default to be english
	}

	//set states
	m.updateThreadsTable()
	m.updateBranchesTable()
	m.updateNotesTable()
	m.updateRecentTable()
	m.updateStatusBar()
	m.DistributeState(&s.App)

	return m
}

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	// Blink ?
	// Start the blinking cursor and the ticker for updating seconds
	return tea.Batch(textinput.Blink, tick())
}

// tick returns a command that sends a tickMsg every second.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NOTE: This only renders the table. Context switching must be done separately via app.UpdateCurrentList()
func (m *Model) updateThreadsTable() {
	var threads []*models.Thread

	threads = m.app.GetThreadList()

	rows := make([]table.Row, len(threads))
	for i, thread := range threads {

		name := thread.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", thread.ID)
		timeStr := thread.CreatedAt.Format("06-01-02 15:04")
		if thread.ID == 0 { // Pending note
			idStr = "P" // Indicate pending
			timeStr = time.Now().Format("06-01-02 15:04")
		}

		bsStrRaw := "0"
		bsStrRaw = strconv.Itoa(len(thread.Branches))

		flagStrRaw := ""
		if thread.Highlight {
			flagStrRaw += "H"
		}
		if thread.Private {
			flagStrRaw += "P"
		}

		// This needs further tuning.
		rows[i] = table.Row{
			idStr,
			timeStr,
			name,
			bsStrRaw,
			flagStrRaw,
		}
	}
	m.threadsTable.SetRows(rows)
}

// NOTE: This only renders the table. Context switching must be done separately via app.UpdateCurrentList()
func (m *Model) updateBranchesTable() {
	var branches []*models.Branch

	branches = m.app.GetActiveBranchList()

	rows := make([]table.Row, len(branches))
	for i, branch := range branches {

		name := branch.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", branch.ID)
		timeStr := branch.CreatedAt.Format("06-01-02 15:04")
		if branch.ID == 0 { // Pending note
			idStr = "P" // Indicate pending
			timeStr = time.Now().Format("06-01-02 15:04")
		}

		nsStrRaw := "0"
		nsStrRaw = strconv.Itoa(len(branch.Notes))

		flagStrRaw := ""
		if branch.Highlight {
			flagStrRaw += "H"
		}
		if branch.Private {
			flagStrRaw += "P"
		}

		// This needs further tuning.
		rows[i] = table.Row{
			idStr,
			timeStr,
			name,
			nsStrRaw,
			flagStrRaw,
		}
	}
	m.branchesTable.SetRows(rows)
}

// NOTE: This only renders the table. Context switching must be done separately via app.UpdateCurrentList()
func (m *Model) updateNotesTable() {
	var selectedNotes []*models.Note
	selectedNotes = m.app.GetActiveNoteList()

	rows := make([]table.Row, len(selectedNotes))
	for i, note := range selectedNotes {

		content := note.Content
		if len(content) > 38 {
			content = content[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", note.ID)
		timeStr := note.CreatedAt.Format("06-01-02 15:04")
		if note.ID == 0 { // Pending note
			idStr = "P" // Indicate pending
			timeStr = time.Now().Format("06-01-02 15:04")
		}

		flagStrRaw := ""
		if note.Highlight {
			flagStrRaw += "H"
		}
		if note.Private {
			flagStrRaw += "P"
		}

		rows[i] = table.Row{
			idStr,
			timeStr,
			content,
			flagStrRaw,
		}
	}
	m.notesTable.SetRows(rows)
}

// This is wrong, to be modified
func (m *Model) updateChangelogTable() {
	editMap := m.app.GetEditMap()
	rows := make([]table.Row, 0, len(editMap))

	for key, edit := range editMap {
		if edit.EditType == -1 { // Skip None edits
			continue
		}

		editTypeName := ""
		switch edit.EditType {
		case 0:
			editTypeName = "Create"
		case 1:
			editTypeName = "Update"
		case 2:
			editTypeName = "Delete"
		case 3:
			editTypeName = "Create"
		case 4:
			editTypeName = "Update"
		case 6:
			editTypeName = "Delete"
		case 7:
			editTypeName = "Create"
		case 8:
			editTypeName = "Update"
		case 10:
			editTypeName = "Delete"
		}

		entityType := key.EntityType
		idStr := fmt.Sprintf("%d", key.ID)
		timeStr := time.Now().Format("06-01-02 15:04")
		description := fmt.Sprintf("%s %s", editTypeName, entityType)

		rows = append(rows, table.Row{
			entityType,
			idStr,
			timeStr,
			description,
		})
	}

	m.changeTable.SetRows(rows)
}

func (m *Model) updateDiffArea(link models.Superlink) {
	// fetch currently active on rencent table
	note := m.app.GetDataMgr().FindNoteByLink(link)
	if note == nil {
		m.diffView.SetContent(fmt.Sprintf(
			"Diff unavailable: note not found with superlink {Thread: %d | Branch: %d | Note: %d}.",
			link.ThreadID,
			link.BranchID,
			link.NoteID,
		))
		return
	}
	m.diffView.SetContent(note.Diff)
	m.diffView.GotoBottom()
}

func (m *Model) updateRecentTable() {
	noteEdits := m.app.GetNoteEditStack()

	rows := make([]table.Row, len(noteEdits))
	for i := range noteEdits {
		noteEdit := noteEdits[len(noteEdits)-1-i]
		link := noteEdit.Link

		// Fetch thread, branch, and note by their IDs
		var threadName, branchName, noteContent string
		flags := ""

		// Get thread name
		if link.ThreadID > 0 {
			thread := m.app.GetDataMgr().FindThreadByID(uint(link.ThreadID))
			if thread != nil {
				threadName = thread.Name
				if len(threadName) > 48 {
					threadName = threadName[:45] + "..."
				}
				if thread.Highlight {
					flags += "TH"
				}
				if thread.Private {
					flags += "TP"
				}
			} else {
				threadName = fmt.Sprintf("T#%d", link.ThreadID)
			}
		} else {
			threadName = "-"
		}

		// Get branch name
		if link.BranchID > 0 {
			branch := m.app.GetDataMgr().FindBranchByID(uint(link.BranchID))
			if branch != nil {
				branchName = branch.Name
				if len(branchName) > 48 {
					branchName = branchName[:45] + "..."
				}
				if branch.Highlight {
					flags += "BH"
				}
				if branch.Private {
					flags += "BP"
				}
			} else {
				branchName = fmt.Sprintf("B#%d", link.BranchID)
			}
		} else {
			branchName = "-"
		}

		// Get note content
		if link.NoteID > 0 {
			note := m.app.GetDataMgr().FindNoteByID(uint(link.NoteID))
			if note != nil {
				noteContent = note.Content
				if len(noteContent) > 68 {
					noteContent = noteContent[:65] + "..."
				}
				if note.Highlight {
					flags += "NH"
				}
				if note.Private {
					flags += "NP"
				}
			} else {
				noteContent = fmt.Sprintf("N#%d", link.NoteID)
			}
		} else {
			noteContent = "-"
		}

		rows[i] = table.Row{
			threadName,
			branchName,
			noteContent,
			flags,
		}
	}
	m.recentTable.SetRows(rows)
}

func (m *Model) printSync() {
	if m.app.Synced {
		m.statusBar.GetTag("Synced").SetValue("Synced")
		m.statusBar.GetTag("Synced").SetColors(colorPtr("232"), colorPtr("118"))
	} else {
		m.statusBar.GetTag("Synced").SetValue("Unsynced")
		m.statusBar.GetTag("Synced").SetColors(colorPtr("232"), colorPtr("208"))
	}
}

// formatTimeAgo returns a human-readable relative time string like "5m ago".
func formatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	d := time.Since(t)
	if d < time.Second*1 {
		return "just now"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds ago", int(d.Minutes()), int(d.Seconds())-60*int(d.Minutes()))
	}
	if d < time.Hour*24 {
		return fmt.Sprintf("%dh%dm ago", int(d.Hours()), int(d.Minutes())-60*int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 7 {
		return fmt.Sprintf("%dd%dh ago", days, int(d.Hours())-24*days)
	}
	// Fallback to date for older items
	// return fmt.Sprintf("%ds ago", int(d))
	return t.Format("01-02 15:04")
}

func (m *Model) updateStatusBar() {
	// Show current table focus
	focusName := "Threads"
	switch m.focus {
	case FocusThreads:
		focusName = "Threads"
		m.statusBar.GetTag("ID").SetValue("#" + strconv.Itoa(int(m.app.GetCurrentThreadID())))
		m.statusBar.GetTag("LastUpdated").SetValue(formatTimeAgo(m.app.GetCurrentThreadLastEdit()))
		m.statusBar.GetTag("Frequency").SetValue(strconv.Itoa(m.app.GetCurrentThreadFrequency()) + " edits")

	case FocusBranches:
		focusName = "Branches"
		m.statusBar.GetTag("ID").SetValue("#" + strconv.Itoa(int(m.app.GetCurrentBranchID())))
		m.statusBar.GetTag("LastUpdated").SetValue(formatTimeAgo(m.app.GetCurrentBranchLastEdit()))
		m.statusBar.GetTag("Frequency").SetValue(strconv.Itoa(m.app.GetCurrentBranchFrequency()) + " edits")

	case FocusNotes:
		focusName = "Notes"
		m.statusBar.GetTag("ID").SetValue("#" + strconv.Itoa(int(m.app.GetCurrentNoteID())))
		m.statusBar.GetTag("LastUpdated").SetValue(formatTimeAgo(m.app.GetCurrentNoteLastEdit()))
		m.statusBar.GetTag("Frequency").SetValue(strconv.Itoa(m.app.GetCurrentNoteFrequency()) + " edits")

	case FocusEdit:
		focusName = "Edit"
		m.statusBar.GetTag("LastUpdated").SetValue("Editing...")
	case FocusChangelog:
		focusName = "Changelog"
	}

	m.statusBar.GetTag("filter").SetValue(focusName)
	m.printSync()
	m.statusBar.GetTag("Time").SetValue(time.Now().Format("15:04:05"))
}
