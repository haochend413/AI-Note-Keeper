package ui

import (
	"fmt"
	"strings"

	"github.com/haochend413/mantis/cmd"
	"github.com/haochend413/mantis/db"
	"github.com/jroimartin/gocui"
)

// quit
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// This sends the current Note buffer to noteDB.
// Resets Buffer.
func sendNote(g *gocui.Gui, view *gocui.View) error {
	//check for empty
	if view.Buffer() == "" {
		return nil
	}
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

// This detects the input command;
// Only the first line of command will be detected and ran
func sendCmd(g *gocui.Gui, view *gocui.View) error {
	//get view content
	if len(view.BufferLines()) == 0 {
		fmt.Println("Invalid")
		return nil
	}
	command := strings.TrimSpace(view.BufferLines()[0])
	cmd.Execute(command)
	return nil
}
