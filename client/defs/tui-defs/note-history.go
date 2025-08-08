package tui_defs

// context in order to identify;
type Context int

const (
	Default Context = iota
	Day
	Week
	Month
	Topic
	Fuzzy
)
