module github.com/haochend413/ntkpr

go 1.25.0

require (
	charm.land/bubbletea/v2 v2.0.6
	github.com/atotto/clipboard v0.1.4
	github.com/haochend413/bubbles/v2 v2.102.0
	github.com/haochend413/lipgloss/v2 v2.100.0
	github.com/sergi/go-diff v1.4.0
	github.com/spf13/cobra v1.10.2
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

// require github.com/charmbracelet/bubbletea v1.3.10

require (
	github.com/asg017/sqlite-vec-go-bindings v0.1.6 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260422141423-a0f1f21775f7 // indirect
	github.com/charmbracelet/x/ansi v0.11.7 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.23 // indirect
	github.com/mattn/go-sqlite3 v1.14.42 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
)

// // These are only for development.
// replace github.com/haochend413/bubbles => ../../bubbles //For development.
// replace github.com/charmbracelet/lipgloss => ../../lipgloss //For development.

replace github.com/charmbracelet/lipgloss v1.1.0 => github.com/haochend413/lipgloss/v2 v0.3.0
replace github.com/haochend413/bubbles v0.2.1 => github.com/haochend413/bubbles/v2 v0.3.0

// replace github.com/haochend413/bubbles/v2 => ../../bubbles //For development.

// replace github.com/haochend413/lipgloss/v2 => ../../lipgloss //For development.
