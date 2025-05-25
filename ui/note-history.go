package ui

import (
	"fmt"

	"github.com/haochend413/mantis/db"
	"github.com/jroimartin/gocui"
)

// delete noteHistory
func quitNoteHistory(g *gocui.Gui, v *gocui.View) error {

	err := g.DeleteView("noteHistory")
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
func setNoteHistory(g *gocui.Gui, view *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("noteHistory", 1, 1, maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note History"
		v.Highlight = true
		v.SelFgColor = gocui.ColorWhite
		// v.Editable = true
		// v.Frame = true
		cursorOn(g, v)

		//every time notehistory is pulled up, it should fetch the newest database note records and show it.
		//also, we should be able to do somthing to the history, not just check, but a way to metitate and reflect. Think!

		//fetch all data from noteDB
		prev := g.CurrentView().Name()
		if _, err := g.SetCurrentView("noteHistory"); err != nil {
			return err
		}
		var history []db.Note
		result := db.NoteDB.Find(&history)

		//display history
		for _, note := range history {
			fmt.Fprintln(v, note.Content)
		}

		if result.Error != nil {
			return result.Error
		}

		noteHistoryKeys(g, prev)
	}

	return nil
}
