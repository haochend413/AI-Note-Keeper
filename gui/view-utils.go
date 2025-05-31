package gui

import (
	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
)

// Control window display
// make sure that view change only happens here for now
// it starts with nothing
var VIEW_SWITCH_HISTORY = []string{""}

// Usually works. Might be buggy
func ToggleWindowDisplay(w *models.Window, g *gocui.Gui) error {
	w.OnDisplay = !w.OnDisplay

	// Safe GUI update
	g.Update(func(g *gocui.Gui) error {
		VIEW_SWITCH_HISTORY = append(VIEW_SWITCH_HISTORY, g.CurrentView().Name())
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
			//manage focus

			if len(VIEW_SWITCH_HISTORY) > 1 {
				VIEW_SWITCH_HISTORY = VIEW_SWITCH_HISTORY[:len(VIEW_SWITCH_HISTORY)-1]
			}
			// Delete/hide the view
			err := g.DeleteView(w.Name)
			g.SetCurrentView(VIEW_SWITCH_HISTORY[len(VIEW_SWITCH_HISTORY)-1])
			if err != nil && err != gocui.ErrUnknownView {
				return err
			}
			w.View = nil
		}
		// fmt.Fprintln(os.Stdout, VIEW_SWITCH_HISTORY[len(VIEW_SWITCH_HISTORY)-1])
		return nil
	})
	return nil
}

// var VIEW_SWITCH_HISTORY []string

// func ToggleWindowDisplay(w *models.Window, g *gocui.Gui) error {
// 	w.OnDisplay = !w.OnDisplay

// 	return g.Update(func(g *gocui.Gui) error {
// 		current := g.CurrentView()
// 		if current != nil && current.Name() != w.Name {
// 			// Only push to history if it's not already the view being toggled
// 			VIEW_SWITCH_HISTORY = append(VIEW_SWITCH_HISTORY, current.Name())
// 		}

// 		if w.OnDisplay {
// 			v, err := g.SetView(w.Name, w.X0, w.Y0, w.X1, w.Y1)
// 			if err != nil && err != gocui.ErrUnknownView {
// 				return err
// 			}
// 			v.Title = w.Title
// 			w.View = v
// 			// Switch focus to the new view
// 			g.SetCurrentView(w.Name)
// 		} else {
// 			// Hide/delete the view
// 			err := g.DeleteView(w.Name)
// 			if err != nil && err != gocui.ErrUnknownView {
// 				return err
// 			}
// 			w.View = nil

// 			// Pop from history and try to set the previous view
// 			for len(VIEW_SWITCH_HISTORY) > 0 {
// 				last := VIEW_SWITCH_HISTORY[len(VIEW_SWITCH_HISTORY)-1]
// 				VIEW_SWITCH_HISTORY = VIEW_SWITCH_HISTORY[:len(VIEW_SWITCH_HISTORY)-1]
// 				if v, _ := g.View(last); v != nil {
// 					g.SetCurrentView(last)
// 					break
// 				}
// 			}
// 		}
// 		return nil
// 	})
// }
