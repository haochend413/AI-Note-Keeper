package gui

import (
	"fmt"
	"strconv"
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
			v.SelBgColor = gocui.ColorBlue
			v.SelFgColor = gocui.ColorYellow
			//here it prints all, and which part gets shown depend on the origin, which we will use to control.
			for _, note := range DB_Data.NoteData {
				var wordlen = 0
				// Remove color codes for length calculation
				timestampRaw := note.CreatedAt.Format("06-01-02 15:04")
				timestamp := "\x1b[35m" + timestampRaw + "\x1b[0m"
				idStr := strconv.FormatUint(uint64(note.ID), 10)
				firstLine := strings.SplitN(note.Content, "\n", 2)[0]

				// Compose the line to print (without color codes for length)
				lineRaw := timestampRaw + "  " + idStr + "  " + firstLine
				if len(firstLine) > 30 {
					lineRaw = timestampRaw + "  " + idStr + "  " + firstLine[:30] + "..."
				}
				wordlen = len(lineRaw)

				// Print with color
				fmt.Fprint(nh, timestamp+"  "+idStr+"  ")
				if len(firstLine) > 30 {
					fmt.Fprint(nh, firstLine[:30])
					fmt.Fprint(nh, " ...")
				} else {
					fmt.Fprint(nh, firstLine)
					if strings.Contains(note.Content, "\n") {
						fmt.Fprint(nh, " ...")
						wordlen += 3
					}
				}

				// Pad spaces to fill the line
				x, _ := nh.Size()
				pad := x - wordlen
				if pad > 0 {
					fmt.Fprint(nh, strings.Repeat(" ", pad))
				}
				fmt.Fprintln(nh, " ")
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
