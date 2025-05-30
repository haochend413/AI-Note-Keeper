package gui

import (
	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
)

// Control window display
func ToggleWindowDisplay(w *models.Window, g *gocui.Gui) error {
	w.OnDisplay = !w.OnDisplay

	// Safe GUI update
	g.Update(func(g *gocui.Gui) error {
		if w.OnDisplay {
			// Create/show the view
			v, err := g.SetView(w.Name, w.X0, w.Y0, w.X1, w.Y1)
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			v.Title = w.Title
			w.View = v
			g.SetCurrentView(w.Name)
		} else {
			// Delete/hide the view
			err := g.DeleteView(w.Name)
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			w.View = nil
		}
		return nil
	})
	return nil
}
