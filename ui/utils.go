package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// quit
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// delete commandInput
func quitCommandInput(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	err := g.DeleteView("commandInput")
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	globalKeys(g)
	return nil
}

// push up input commandbar, setview
func setCommandInput(g *gocui.Gui, view *gocui.View) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("commandInput", 20, maxY/2-10, maxX, maxY/2+10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command Input"
		v.Editable = true
		v.Frame = true
		fmt.Fprintln(v, ":")
	}

	if _, err := g.SetCurrentView("commandInput"); err != nil {
		return err
	}
	cmdinputKeys(g)

	return nil
}
