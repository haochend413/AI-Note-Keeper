package gui

import (
	"github.com/haochend413/mantis/controllers"
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
		g.CurrentView().FrameColor = gocui.ColorGreen
		return nil
	case "note-history":
		g.CurrentView().FrameColor = gocui.ColorWhite

		g.SetCurrentView("note")
		g.CurrentView().FrameColor = gocui.ColorGreen
		return nil
	default:
		return nil
	}
}

// // Move Cursor; Update View Related Values: gui imports views
// func (gui *Gui) HandleNoteHistoryCursorMove(direction string) func(*gocui.Gui, *gocui.View) error {
// 	return func(g *gocui.Gui, v *gocui.View) error {
// 		switch direction {
// 		//move cursor up
// 		case "up":
// 			err := controllers.CursorUp(gui.windows[1].View)
// 			_, views.P_CURSOR_NH = gui.windows[1].View.Cursor()
// 			views.UpdateSelectedNote(gui.g, DB_Data)
// 			// fmt.Fprintln(os.Stdout, P_CURSOR_NH)
// 			return err
// 		//move cursor down
// 		case "down":
// 			err := controllers.CursorDown(gui.windows[1].View)
// 			_, views.P_CURSOR_NH = gui.windows[1].View.Cursor()
// 			views.UpdateSelectedNote(gui.g, DB_Data)
// 			// fmt.Fprintln(os.Stdout, P_CURSOR_NH)
// 			return err
// 		case "left":
// 			return controllers.CursorLeft(gui.windows[1].View)
// 		case "right":
// 			return controllers.CursorRight(gui.windows[1].View)
// 		default:
// 			return nil
// 		}
// 	}
// }

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
			}
			if views.P_ORIGIN_NH > 0 {
				views.P_ORIGIN_NH -= 1
			}

			//re-render
			views.UpdateSelectedNote(gui.g, DB_Data)
			return views.UpdateHistoryDisplay(v)
			// return views.UpdateHistoryDisplay(v)
		//move cursor down
		case "down":
			// two things to note;
			_, height := v.Size()

			// Need to know total lines in view content to avoid moving beyond content
			lines := len(v.BufferLines())

			if views.P_ORIGIN_NH+views.P_CURSOR_NH+1 < lines {
				if views.P_CURSOR_NH < height-1 {
					views.P_CURSOR_NH += 1
					if err := v.SetCursor(0, views.P_CURSOR_NH+1); err != nil {
						views.P_CURSOR_NH -= 1
						views.P_ORIGIN_NH += 1
						return v.SetOrigin(0, views.P_ORIGIN_NH+1)
					}
				} else {
					views.P_ORIGIN_NH += 1
					return v.SetOrigin(0, views.P_ORIGIN_NH+1)
				}
			}
			// controllers.CursorDown(gui.windows[1].View)
			// if views.P_CURSOR_NH == winheight-1 {
			// 	dlen := len(DB_Data.NoteDBData)
			// 	views.OriginDown(winheight, dlen)
			// } else {
			// 	if views.P_CURSOR_NH < len(DB_Data.NoteDBData)-1 {
			// 		views.P_CURSOR_NH += 1
			// 	}
			// }
			// //re-render
			// // fmt.Fprintln(gui.windows[0].View, views.P_CURSOR_NH)
			// // _, y := v.Size()
			// // fmt.Fprintln(gui.windows[0].View, y)
			// views.UpdateSelectedNote(gui.g, DB_Data)
			// v.SetOrigin(0, views.P_ORIGIN_NH)
			// // fmt.Fprintln(os.Stdout, "test1")

			// //cursor
			// cx, _ := v.Cursor()
			// v.SetCursor(cx, views.P_CURSOR_NH)
			// // views.UpdateHistoryDisplay(v)
			// // a, _ := v.Cursor()
			// // fmt.Fprintln(gui.windows[0].View, views.P_CURSOR_NH)
			// // v.MoveCursor(0, 1)
			// // return nil
			views.UpdateSelectedNote(gui.g, DB_Data)
			return views.UpdateHistoryDisplay(v)
		case "left":
			return controllers.CursorLeft(gui.windows[1].View)
		case "right":
			return controllers.CursorRight(gui.windows[1].View)
		default:
			return nil
		}
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
	return nil
}

// func (gui *Gui) HandleSwitchLine(g *gocui.Gui, v *gocui.View) error {
// 	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
// 	v.edi()
// 	return nil
// }
