package gui

import (
	"github.com/haochend413/mantis/gui/controllers"
	"github.com/jroimartin/gocui"
)

// Define layout for all views;
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("note", 1, maxY/5*4, maxX-1, maxY/5*4+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note"
		controllers.CursorOn(g, v)
	}
	if v, err := g.SetView("note-history", 1, 1, maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note History"
		controllers.CursorOn(g, v)
	}
	if v, err := g.SetView("cmd", 20, maxY/2-1, maxX-20, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Cmd"
		controllers.CursorOn(g, v)
	}
	g.SetCurrentView("cmd")
	return nil
}
