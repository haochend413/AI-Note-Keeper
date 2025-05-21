package ui

import (
	"strings"

	"github.com/haochend413/mantis/db"
	"github.com/jroimartin/gocui"
)

// quit
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// cursor configs

func cursorOn(g *gocui.Gui, view *gocui.View) error {
	g.Cursor = true
	lines := view.BufferLines()

	// // Remove trailing empty lines
	// for len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
	// 	lines = lines[:len(lines)-1]
	// }

	if len(lines) == 0 {
		return view.SetCursor(0, 0)
	}

	px := len(lines[len(lines)-1])
	py := len(lines) - 1
	return view.SetCursor(px, py)
}

// func cursorOff(g *gocui.Gui, view *gocui.View) {
// 	g.Cursor = false
// }

// note submit: to database
func sendNote(g *gocui.Gui, view *gocui.View) error {
	s := strings.TrimSpace(view.Buffer())

	if err := db.AddNote(s); err != nil {
		return err
	}

	//clear note view
	view.Clear()
	// reset cursor
	cursorOn(g, view)

	return nil
}
