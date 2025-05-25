package ui

import (
	"github.com/jroimartin/gocui"
)

// delete commandInput
func quitCommandInput(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	//delete cursor
	// cursorOff(g, v)

	err := g.DeleteView("commandInput")
	g.Cursor = false
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	prev := g.CurrentView().Name()
	g.SetCurrentView("note")
	cursorOn(g, g.CurrentView())
	noteKeys(g, prev)
	return nil
}

// push up input commandbar, setview
func setCommandInput(g *gocui.Gui, view *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("commandInput", 20, maxY/2-1, maxX-20, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command Input"
		v.Editable = true
		v.Frame = true

		// //open cursor
		// x, y := v.Cursor()
		// v.SetCursor(x+1, y) // move right 1 colum
		cursorOn(g, v)

	}

	prev := g.CurrentView().Name()

	// fmt.Println(g.CurrentView())
	if _, err := g.SetCurrentView("commandInput"); err != nil {
		return err
	}
	// fmt.Println(g.CurrentView())

	cmdinputKeys(g, prev)

	return nil
}
