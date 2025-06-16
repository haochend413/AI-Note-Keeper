package keybindings

import "github.com/charmbracelet/bubbles/key"

// model
type GlobalKeyMap struct {
	QuitApp key.Binding
}

// init
var GlobalKeys = GlobalKeyMap{
	QuitApp: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
