package gui

import (
	"log"

	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
)

// main Gui struct
type Gui struct {
	g       *gocui.Gui
	windows []*models.Window
}

// need to use a map to hande quick window search

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
	// Set layout manager function (called every frame to layout views)
	gui.g.SetManagerFunc(gui.layout)
	gui.windows = gui.CreateWindowTemplates()

	//init keybindings
	if err := gui.InitKeyBindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
