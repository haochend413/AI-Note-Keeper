package gui

import (
	"github.com/jroimartin/gocui"
)

var GUI *gocui.Gui

func AppInit() {
	gui := Gui{}
	//start up my gui
	gui.GuiInit()
}
