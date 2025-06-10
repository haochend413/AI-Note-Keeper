package gui

import (
	"github.com/awesome-gocui/gocui"
)

// var GUI *gocui.Gui

func AppInit() {
	gui := &Gui{}
	gui.GuiInit()
	// defer gui.GuiClose()
}

func (gui *Gui) AppQuit() error {
	gui.GuiClose()
	return gocui.ErrQuit
}
