package views

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/haochend413/mantis/models"
)

var VIEW_SWITCH_HISTORY = []string{""}

// Switch window on/off and change focus.
// Usually works. Might be buggy
func ToggleWindowDisplay(w *models.Window, g *gocui.Gui) error {
	w.OnDisplay = !w.OnDisplay

	// Safe GUI update
	g.Update(func(g *gocui.Gui) error {
		VIEW_SWITCH_HISTORY = append(VIEW_SWITCH_HISTORY, g.CurrentView().Name())
		if w.OnDisplay {
			// Create/show the view
			v, err := g.SetView(w.Name, w.X0, w.Y0, w.X1, w.Y1, 0)
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
		return nil
	})
	return nil
}

/*
List Display
*/

// Note History Display
var P_ORIGIN_NH int // Position of origin (start line inside data)
var P_CURSOR_NH int

// Position of cursor (relative to window)

// Update Window Display based on what we have;
func UpdateHistoryDisplay(v *gocui.View) error {
	//display content
	v.SetOrigin(0, P_ORIGIN_NH)
	// fmt.Fprintln(os.Stdout, "test1")
	//cursor
	cx, _ := v.Cursor()
	v.SetCursor(cx, P_CURSOR_NH)

	// fmt.Fprintln(os.Stdout, "test2")
	return nil
	// v.MoveCursor(0, 1)
	// return nil
}

// set origin
func OriginDown(windowHeight int, dataLength int) {
	if P_ORIGIN_NH < dataLength-windowHeight {
		P_ORIGIN_NH += 1
	}
}

func OriginUp() {
	if P_ORIGIN_NH > 0 {
		P_ORIGIN_NH -= 1
	}
}

// This will display the info of the note at P_CURSOR_NH on note-details
func UpdateSelectedNote(g *gocui.Gui, data *models.DB_Data) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("note-detail")
		if err != nil {
			return nil
		}
		v.Clear()
		v.Wrap = true
		fmt.Fprint(v, data.NoteDBData[P_CURSOR_NH+P_ORIGIN_NH-1].Content)
		return nil
	})
	return nil
}
