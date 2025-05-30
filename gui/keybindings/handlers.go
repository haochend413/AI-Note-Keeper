package keybindings

import (
	"github.com/haochend413/mantis/gui/controllers"
	"github.com/jroimartin/gocui"
)

func HandleQuit(g *gocui.Gui, v *gocui.View) error {
	return controllers.QuitApp()
}
