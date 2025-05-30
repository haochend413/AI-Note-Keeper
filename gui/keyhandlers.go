package gui

import (
	"github.com/haochend413/mantis/controllers"
	"github.com/jroimartin/gocui"
)

func HandleAppQuit(g *gocui.Gui, v *gocui.View) error {
	return controllers.QuitApp()

}

// View setup;
func (gui *Gui) HandleNoteDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return ToggleWindowDisplay(gui.windows[0], gui.g)
}

// View setup;
func (gui *Gui) HandleNoteHistoryDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return ToggleWindowDisplay(gui.windows[1], gui.g)
}

// View setup;
func (gui *Gui) HandleCmdDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return ToggleWindowDisplay(gui.windows[2], gui.g)
}
