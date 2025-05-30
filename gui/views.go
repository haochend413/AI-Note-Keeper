package gui

import (
	"github.com/haochend413/mantis/gui/models"
	"github.com/jroimartin/gocui"
)

// This defines all the views existing in the main gui
type Windows struct {
	Note        *models.Window
	NoteHistory *models.Window
	Cmd         *models.Window
}

type windowNameMapping struct {
	window *models.Window
	name   string
}

// Use pointer: manage the real state;
func (gui *Gui) MapWindowNames() []*windowNameMapping {
	return []*windowNameMapping{
		{window: gui.windows.Note, name: "note"},
		{window: gui.windows.NoteHistory, name: "note-history"},
		{window: gui.windows.Cmd, name: "cmd"},
	}
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10)
}

// create windows and include configs
func (gui *Gui) createAllWindows() error {
	for _, w := range gui.MapWindowNames() {
		//here it set up view prepare;
		v, err := gui.prepareView(w.name)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		//init
		w.window = &models.Window{}
		w.window.View = v
		//This implicitly sets that the name of the view is the same as the name of the Window;
		w.window.Name = w.name
		w.window.OnDisplay = true

	}
	return nil
}
