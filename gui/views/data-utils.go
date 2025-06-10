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
	data.NoteDBData = append(data.NoteDBData, note)
	// reset Note view
	g.CurrentView().Clear()
	controllers.CursorOn(g, g.CurrentView())
	return nil
}

/*
Note-Detail
*/
