package keybindings

import "github.com/charmbracelet/bubbles/key"

// model
type GlobalKeyMap struct {
	QuitApp     key.Binding
	SwitchFocus key.Binding
}

type DailyKeyMap struct {
	ToggleSuccess key.Binding
}

// init
var GlobalKeys = GlobalKeyMap{
	QuitApp: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	// SwitchFocus: key.NewBinding(
	// 	key.WithKeys("tab"),
	// 	key.WithHelp("tab", "switch focus"),
	// ),
}

var DailyKeys = DailyKeyMap{
	ToggleSuccess: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "toggle success"),
	),
}
