package keybindings

import "github.com/charmbracelet/bubbles/key"

// model
type GlobalKeyMap struct {
	QuitApp     key.Binding
	SwitchFocus key.Binding
}

type NoteKeyMap struct {
	SendNote key.Binding
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

var Notekeys = NoteKeyMap{
	SendNote: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send note to database"),
	),
}
