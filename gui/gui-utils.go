package gui

import "github.com/jroimartin/gocui"

// Maybe a bad idea;
func (gui *Gui) G() *gocui.Gui {
	return gui.g
}
