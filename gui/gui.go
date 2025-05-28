package gui

import (
	"log"

	"github.com/jroimartin/gocui"
)

// main Gui struct
type Gui struct {
	g     *gocui.Gui
	views Views
}

// Define layout for all views;
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("note", 1, maxY/5*4, maxX-1, maxY/5*4+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note"
		// CursorOn(g, v)
	}
	if v, err := g.SetView("note-history", 1, 1, maxX/3, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Note History"
		// CursorOn(g, v)
	}
	if v, err := g.SetView("cmd", 20, maxY/2-1, maxX-20, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Cmd"
		// CursorOn(g, v)
	}

	return nil
}

// This function inits a new Gui object;
func (gui *Gui) GuiInit() {
	// setup the new gui instance
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		//check startup err
		log.Panicln(err)
	}
	gui.g = g
	defer g.Close()
	//
	gui.createAllViews() //create all the views;
	// Set layout manager function (called every frame to layout views)
	gui.g.SetManagerFunc(layout)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// This function quits gui
func (gui *Gui) QuitGui() error {
	return gocui.ErrQuit
}
