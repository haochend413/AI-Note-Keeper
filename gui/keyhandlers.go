package gui

import (
	"github.com/haochend413/mantis/gui/views"

	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) HandleAppQuit(g *gocui.Gui, v *gocui.View) error {
	return gui.AppQuit()

}

func (gui *Gui) HandleDataUpdate(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprintln(v, "test")
	return gui.DBManager.RefreshAll(DB_Data)
	// return nil
}

// View setup;
func (gui *Gui) HandleNoteDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return views.ToggleWindowDisplay(gui.windows[0], gui.g)
}

// View setup;
func (gui *Gui) HandleNoteHistoryDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return views.ToggleWindowDisplay(gui.windows[1], gui.g)
}

// View setup;
func (gui *Gui) HandleCmdDisplay(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	return views.ToggleWindowDisplay(gui.windows[2], gui.g)
}

// View switch
// Should not go to read-only views
func (gui *Gui) HandleViewLoop(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "note":
		g.CurrentView().FrameColor = gocui.ColorWhite
		g.SetCurrentView("note-history")
		g.Cursor = false
		g.CurrentView().FrameColor = gocui.ColorGreen
		return nil
	case "note-history":
		g.CurrentView().FrameColor = gocui.ColorWhite

		g.SetCurrentView("note")
		g.Cursor = true
		g.CurrentView().FrameColor = gocui.ColorGreen
		return nil
	default:
		return nil
	}
}

/*
Note History Display
*/
// Move Cursor; Update View Related Values: gui imports views
func (gui *Gui) HandleHistorySelect(direction string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		// _, winheight := v.Size()
		switch direction {
		//move cursor up
		case "up":
			// if  at top, move origin
			if views.P_CURSOR_NH > 0 {
				views.P_CURSOR_NH -= 1
			} else {
				if views.P_ORIGIN_NH > 0 {
					views.P_ORIGIN_NH -= 1
				}
			}

			//re-render
			views.UpdateSelectedNote(gui.g, DB_Data)
			// return views.UpdateHistoryDisplay(v)
			return nil

			// return views.UpdateHistoryDisplay(v)
		//move cursor down
		case "down":
			// two things to note;
			_, height := v.Size()

			// Need to know total lines in view content to avoid moving beyond content
			lines := len(DB_Data.NoteDBData)

			//if there is need to move cursor down
			if views.P_ORIGIN_NH+views.P_CURSOR_NH < lines-1 {
				if views.P_CURSOR_NH < height-1 {
					views.P_CURSOR_NH += 1
				} else {
					// if cursor is at the bottom, move origin
					views.P_ORIGIN_NH += 1
					// views.P_CURSOR_NH += 1
				}
			}
			views.UpdateSelectedNote(gui.g, DB_Data)
			// return views.UpdateHistoryDisplay(v)
			return nil
		case "left":
			//move up 5
			// if  at top, move origin

			//calculate separately
			D_Origin := min(0, views.P_CURSOR_NH-5)
			if D_Origin < 0 {
				//reset cursor
				views.P_CURSOR_NH = 0
				//move origin
				views.P_ORIGIN_NH += D_Origin
				//check for upper bound
				if views.P_ORIGIN_NH < 0 {
					views.P_ORIGIN_NH = 0
				}
			} else {
				//just move cursor
				views.P_CURSOR_NH -= 5
			}

			//re-render
			views.UpdateSelectedNote(gui.g, DB_Data)
			// return views.UpdateHistoryDisplay(v)
			return nil
		case "right":
			//move down 5
			// two things to note;
			_, height := v.Size()
			// Need to know total lines in view content to avoid moving beyond content

			D_Origin := max(0, views.P_CURSOR_NH+5-height)
			if D_Origin > 0 {
				//reset cursor
				views.P_CURSOR_NH = height - 1
				//move origin
				views.P_ORIGIN_NH += D_Origin
				//check for upper bound
				if views.P_ORIGIN_NH > len(DB_Data.NoteDBData)-height {
					views.P_ORIGIN_NH = len(DB_Data.NoteDBData) - height
				}
			} else {
				//just move cursor
				views.P_CURSOR_NH += 5
				if views.P_CURSOR_NH >= height {
					views.P_CURSOR_NH = height - 1
				}
			}
			views.UpdateSelectedNote(gui.g, DB_Data)
			// return views.UpdateHistoryDisplay(v)
			return nil
		default:
			return nil
		}
	}
}

func (gui *Gui) HandleJumpToEnd() func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		//jump to end
		_, height := v.Size()
		if len(DB_Data.NoteDBData) < height {
			views.P_CURSOR_NH = len(DB_Data.NoteDBData) - 1
			views.P_ORIGIN_NH = 0
			// return nil
		} else {
			views.P_CURSOR_NH = height - 1
			views.P_ORIGIN_NH = len(DB_Data.NoteDBData) - height
			// return nil
		}
		views.UpdateSelectedNote(gui.g, DB_Data)
		return nil
	}
}

/*
Note view
*/

// Send Note.
// Update history & detail demo.
func (gui *Gui) HandleSendNote(g *gocui.Gui, v *gocui.View) error {
	// update db data
	views.SendNote(gui.windows[0], gui.g, DB_Data)
	_, y := gui.windows[1].View.Size()
	views.P_ORIGIN_NH = max(0, len(DB_Data.NoteDBData)-y)

	// move cursor to last visible line
	views.P_CURSOR_NH = min(len(DB_Data.NoteDBData)-1, y-1)

	views.UpdateHistoryDisplay(gui.windows[1].View)

	// demo selected note
	return views.UpdateSelectedNote(gui.g, DB_Data)
}

// func (gui *Gui) HandleSwitchLine(g *gocui.Gui, v *gocui.View) error {
// 	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
// 	v.edi()
// 	return nil
// }
