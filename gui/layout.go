package gui

import (
	"github.com/jroimartin/gocui"
)

var FIRSTINITCHECK bool = true

// Define layout for all views;
func (gui *Gui) layout(g *gocui.Gui) error {
	//init template
	if FIRSTINITCHECK {
		gui.windows = gui.CreateWindowTemplates()
		FIRSTINITCHECK = false
	}

	//here, only check logic
	// init views
	for _, w := range gui.windows {
		// fmt.Fprint(os.Stdout, w.Name)
		if !w.OnDisplay {
			// Don't show views that are off
			continue
		}
		//here it set up view prepare;
		// Only initialize if the view was just created

		v, err := gui.prepareView(w)
		if !w.OnDisplay {
			g.DeleteView(w.Name)
		}
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}

		//check: init if only first created
		if err == gocui.ErrUnknownView {
			v.Title = w.Title
			w.View = v
			// Optional: set cursor, wrap, etc.
			// v.Wrap = true
		}

	}

	//setstartview
	g.SetCurrentView("note")

	return nil
}
