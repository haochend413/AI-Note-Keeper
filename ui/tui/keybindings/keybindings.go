package keybindings

import "github.com/charmbracelet/bubbles/key"

// model
type GlobalKeyMap struct {
	QuitApp     key.Binding
	SwitchFocus key.Binding
}

type NoteKeyMap struct {
	SendNote  key.Binding
	SendTopic key.Binding
}

type HistoryKeyMap struct {
	DayContext     key.Binding
	WeekContext    key.Binding
	MonthContext   key.Binding
	DefaultContext key.Binding
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
}

// history
var Historykeys = HistoryKeyMap{
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
}
