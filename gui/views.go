package gui

import "github.com/jroimartin/gocui"

// This defines all the views existing in the main gui
type Views struct {
	Note        *gocui.View
	NoteHistory *gocui.View
	Cmd         *gocui.View
}

type viewNameMapping struct {
	view *gocui.View
	name string
}

// Use pointer: manage the real state;
func (gui *Gui) MapViewNames() []*viewNameMapping {
	return []*viewNameMapping{
		{view: gui.views.Note, name: "note"},
		{view: gui.views.NoteHistory, name: "note-history"},
		{view: gui.views.Cmd, name: "cmd"},
	}
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10)
}

// create views and include configs
func (gui *Gui) createAllViews() error {
	for _, vm := range gui.MapViewNames() {
		//here it set up view prepare;
		v, err := gui.prepareView(vm.name)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		vm.view = v
	}
	return nil
}
