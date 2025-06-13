package views

import (
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/haochend413/mantis/controllers"
	"github.com/haochend413/mantis/models"
)

// Control window display
// make sure that view change only happens here for now
// it starts with nothing

// fetch the current content input of that view;
func FetchContent(w *models.Window, g *gocui.Gui) string {
	return strings.TrimSpace(w.View.Buffer())
}

/*
Note View
*/
// store current note to DB_Data
func SendNote(w *models.Window, g *gocui.Gui, data *models.DB_Data) error {
	content := FetchContent(w, g)
	if content == "" {
		return nil
	}
	note := &models.Note{Content: content}
	data.NoteData = append(data.NoteData, note)
	// reset Note view
	g.CurrentView().Clear()
	controllers.CursorOn(g, g.CurrentView())
	return nil
}

/*
Note History View
*/
// remove the note at the current index: P_CURSOR_NH + P_ORIGIN_NH
func DeleteNote(w *models.Window, g *gocui.Gui, data *models.DB_Data) error {
	//check for boundary
	//if no notes, do nothing
	_, height := w.View.Size()
	if len(data.NoteData) == 0 {
		return nil
	}

	// delete current note
	data.NoteData = append(data.NoteData[:P_CURSOR_NH+P_ORIGIN_NH], data.NoteData[P_CURSOR_NH+P_ORIGIN_NH+1:]...)

	//aftter delete, check for valid position, if not valid, return to last valid
	if P_CURSOR_NH+P_ORIGIN_NH >= len(data.NoteData) {
		if P_CURSOR_NH == height-1 {
			//at bottom, move origin up
			P_ORIGIN_NH = max(0, P_ORIGIN_NH-1)
		} else {
			//just move cursor up
			P_CURSOR_NH = max(0, P_CURSOR_NH-1)
		}
		return nil
	}
	UpdateSelectedNote(g, data)
	return nil

}

/*
Note-Detail
*/
