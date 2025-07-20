package keybindings

import "github.com/charmbracelet/bubbles/key"

// model
type GlobalKeyMap struct {
	QuitApp     key.Binding
	SwitchFocus key.Binding
	FetchFromDB key.Binding
}

type NoteKeyMap struct {
	SendNote       key.Binding
	SendTopic      key.Binding
	ToggleEditable key.Binding
}

type HistoryKeyMap struct {
	DayContext     key.Binding
	WeekContext    key.Binding
	MonthContext   key.Binding
	DefaultContext key.Binding
	DeleteNote     key.Binding
	EditNote       key.Binding
}

// init
var GlobalKeys = GlobalKeyMap{
	QuitApp: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	FetchFromDB: key.NewBinding(
		key.WithKeys("ctrl+q"),
		key.WithHelp("ctrl+q", "switch focus"),
	),
}

// note
var Notekeys = NoteKeyMap{
	SendNote: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send note to database"),
	),
	SendTopic: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add topic to database"),
	),
	ToggleEditable: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "toggle edit"),
	),
}

// history
var Historykeys = HistoryKeyMap{
	//Context switch
	DayContext: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Day view"),
	),
	WeekContext: key.NewBinding(
		key.WithKeys("W"),
		key.WithHelp("W", "Week view"),
	),
	MonthContext: key.NewBinding(
		key.WithKeys("M"),
		key.WithHelp("M", "Month view"),
	),
	DefaultContext: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "Full view"),
	),
	//Functionalities
	DeleteNote: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "Delete Selected Note"),
	),
	EditNote: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Edit Current Note"),
	),
}
