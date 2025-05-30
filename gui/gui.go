package gui

import (
	"log"

	"github.com/jroimartin/gocui"
)

// main Gui struct
type Gui struct {
	g       *gocui.Gui
	windows Windows
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
	defer gui.g.Close()
	//
	gui.createAllWindows() //create all the views

	// Set layout manager function (called every frame to layout views)
	gui.g.SetManagerFunc(layout)

	//init keybindings
	if err := gui.InitKeyBindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
