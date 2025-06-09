package gui

import (
	"github.com/haochend413/mantis/controllers"
	"github.com/haochend413/mantis/gui/views"

	"github.com/jroimartin/gocui"
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
		g.SetCurrentView("note-history")
		return nil
	case "note-history":
		g.SetCurrentView("note")
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
// 			_, views.CURRENT_NOTE_INDEX = gui.windows[1].View.Cursor()
// 			views.UpdateSelectedNote(gui.g, DB_Data)
// 			// fmt.Fprintln(os.Stdout, Current_Note_Index)
// 			return err
// 		//move cursor down
// 		case "down":
// 			err := controllers.CursorDown(gui.windows[1].View)
// 			_, views.CURRENT_NOTE_INDEX = gui.windows[1].View.Cursor()
// 			views.UpdateSelectedNote(gui.g, DB_Data)
// 			// fmt.Fprintln(os.Stdout, Current_Note_Index)
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
		_, winheight := v.Size()
		switch direction {
		//move cursor up
		case "up":
			// if  at top, move origin
			if views.CURRENT_NOTE_INDEX == 0 {
				views.OriginUp()
			} else {
				views.CURRENT_NOTE_INDEX -= 1
			}
			//re-render
			views.UpdateSelectedNote(gui.g, DB_Data)
			return views.UpdateHistoryDisplay(v)
		//move cursor down
		case "down":
			if views.CURRENT_NOTE_INDEX == winheight-1 {
				dlen := len(DB_Data.NoteDBData)
				views.OriginDown(winheight, dlen)
			} else {
				if views.CURRENT_NOTE_INDEX < len(DB_Data.NoteDBData)-1 {
					views.CURRENT_NOTE_INDEX += 1
				}
			}
			//re-render
			// fmt.Fprintln(gui.windows[0].View, views.CURRENT_NOTE_INDEX)
			// _, y := v.Size()
			// fmt.Fprintln(gui.windows[0].View, y)
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
	views.CURRENT_NOTE_INDEX = min(len(DB_Data.NoteDBData)-1, y-1)

	views.UpdateHistoryDisplay(gui.windows[1].View)

	// demo selected note
	return views.UpdateSelectedNote(gui.g, DB_Data)
}

func (gui *Gui) HandleSwitchLine(g *gocui.Gui, v *gocui.View) error {
	// fmt.Fprint(os.Stdout, gui.windows[0].Name)
	v.EditNewLine()
	return nil
}
