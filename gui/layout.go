package gui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/haochend413/mantis/controllers"
	"github.com/haochend413/mantis/gui/views"
)

// Define layout for all views;
func (gui *Gui) layout(g *gocui.Gui) error {
	//init template
	if gui.first_init_check {
		gui.windows = gui.CreateWindowTemplates()
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
		// //Dont know why here, but might be useful
		// if !w.OnDisplay {
		// 	g.DeleteView(w.Name)
		// }
		// v.BgColor = gocui.ColorGreen
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}

		//check: init if only first created
		if err == gocui.ErrUnknownView {
			//view config
			v.Title = w.Title
			w.View = v
			if w.Editable {
				v.Editable = true
			}
			if w.Cursor {
				controllers.CursorOn(g, v)
			}
		}

		//view-specific logic here
		//fetch again
		if w.Name == "note-history" {
			// g.Cursor = false
			nh := w.View
			nh.Clear()
			//display history
			nh.Highlight = true
			v.SelBgColor = gocui.ColorCyan
			v.SelFgColor = gocui.ColorBlue
			//here it prints all, and which part gets shown depend on the origin, which we will use to control.
			for _, note := range DB_Data.NoteDBData {
				timestamp := "\x1b[35m" + note.CreatedAt.Format("06-01-02 15:04") + "\x1b[0m"
				fmt.Fprint(nh, timestamp)
				fmt.Fprint(nh, "  ")
				fmt.Fprint(nh, note.ID)
				fmt.Fprint(nh, "  ")
				firstLine := strings.SplitN(note.Content, "\n", 2)[0]
				var d = false
				if len(strings.SplitN(note.Content, "\n", 2)) > 1 {
					d = true
				}

				if d {
					if len(firstLine) <= 30 {
						fmt.Fprint(nh, firstLine)
						fmt.Fprintln(nh, "...")
					} else {
						fmt.Fprint(nh, firstLine[:30])
						fmt.Fprintln(nh, "...")
					}
				} else {
					if len(firstLine) <= 30 {
						fmt.Fprintln(nh, firstLine)
					} else {
						fmt.Fprint(nh, firstLine[:30])
						fmt.Fprintln(nh, "...")
					}
				}
			}
			//reset origin & cursor
			v.SetOrigin(0, views.P_ORIGIN_NH)
			v.SetCursor(0, views.P_CURSOR_NH)
		}

		// if w.Name == "note-detail" {
		// 	nh, e := g.View("note-detail")
		// 	nh.Clear()
		// 	// fmt.Fprint(os.Stdout, "hihi, \n hihi, \n hihi")
		// }

	}

	//setstartview
	if gui.first_init_check {
		g.SetCurrentView("note")
		g.CurrentView().FrameColor = gocui.ColorGreen
		gui.first_init_check = false
	}

	return nil
}

// // help get the part of history demonstrated in the note-history view
// func getHistoryDisplay(v *gocui.View) []*models.Note {
// 	_, maxY := v.Size()
// 	if len(DB_Data.NoteDBData) < maxY {
// 		return DB_Data.NoteDBData
// 	}

// 	//if not, return the last portion;
// 	front := max(Current_Note_Index-maxY, 0)
// 	end := min(front+maxY, len(DB_Data.NoteDBData))
// 	return DB_Data.NoteDBData[front:end]
// }
