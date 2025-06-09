package views

import (
	"fmt"

	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
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
		return nil
	})
	return nil
}

/*
List Display
*/

// Note History Display
var P_ORIGIN_NH int // Position of origin (start line inside data)
var CURRENT_NOTE_INDEX int

// Position of cursor (relative to window)

// Update Window Display based on what we have;
func UpdateHistoryDisplay(v *gocui.View) error {
	//display content
	v.SetOrigin(0, P_ORIGIN_NH)
	//cursor
	cx, _ := v.Cursor()
	return v.SetCursor(cx, CURRENT_NOTE_INDEX)
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

// This will display the info of the note at Current_Note_Index on note-details
func UpdateSelectedNote(g *gocui.Gui, data *models.DB_Data) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("note-detail")
		if err != nil {
			return nil
		}
		v.Clear()
		v.Wrap = true
		fmt.Fprint(v, data.NoteDBData[CURRENT_NOTE_INDEX+P_ORIGIN_NH].Content)
		return nil
	})
	return nil
}

// // Calculate the part of notedb data that is to be displayed inside the note-history view and avoid invalid point
// func SelectHistoryDisplay(windowHeight int, data *models.DB_Data) (start int, end int) {
// 	//here to be cautious: just say smaller than
// 	if len(data.NoteDBData) < windowHeight {
// 		return 0, len(data.NoteDBData)
// 	}

// 	//else, we need to check cursor position, which is current note selected
// 	// we also know that the last line lies at windowHeight - 1;
// 	// here we need to make sure that origin and cursor pos stay in there desired range
// 	start = P_ORIGIN_NH + CURRENT_NOTE_INDEX
// 	end = start + windowHeight
// 	return start, end
// }
