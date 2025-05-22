package ui

import (
	"github.com/jroimartin/gocui"
)

// delete commandInput
// Temporarily useless since I do not want to quit note: it should be always there.
// func quitNote(g *gocui.Gui, v *gocui.View) error {
// 	v.Clear()
// 	//delete cursor
// 	// cursorOff(g, v)

// 	err := g.DeleteView("note")
// 	g.Cursor = false
// 	if err != nil && err != gocui.ErrUnknownView {
// 		return err
// 	}
// 	globalKeys(g)
// 	return nil
// }

// push up input commandbar, setview
func setNote(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("note", 1, maxY/4*3, maxX-1, maxY/4*3+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note"
		v.Editable = true
		v.Frame = true

		// //open cursor
		// x, y := v.Cursor()
		// v.SetCursor(x+1, y) // move right 1 colum
		cursorOn(g, v)
	}

	if _, err := g.SetCurrentView("note"); err != nil {
		return err
	}

	//key config
	noteKeys(g)

	return nil
}
