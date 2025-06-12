package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/haochend413/mantis/gui/keybindings"
	"github.com/haochend413/mantis/models"
)

// This function inits all the keybindings
func (gui *Gui) InitKeyBindings() error {
	for _, k := range CreateAllKeybinders(gui) {
		//use parsor function to get key
		s := keybindings.Parsor(k.Key)
		if s.Valid {
			// kCopy := k
			if s.IsRune {
				if err := gui.g.SetKeybinding(k.ViewName, s.Rune, k.Modifier, k.Handler); err != nil {
					return err
				}
			} else {
				if err := gui.g.SetKeybinding(k.ViewName, s.Key, k.Modifier, k.Handler); err != nil {
					return err
				}

			}
		}
	}
	return nil
}

func CreateAllKeybinders(gui *Gui) []*models.KeyBinder {
	return []*models.KeyBinder{
		{
			ViewName: "",
			Key:      "ct-c",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleAppQuit,
		},
		{
			ViewName: "",
			Key:      "ct-a",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleNoteDisplay,
		},
		{
			ViewName: "",
			Key:      "ct-e",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleNoteHistoryDisplay,
		},
		{
			ViewName: "",
			Key:      "ct-x",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleCmdDisplay,
		},
		{
			ViewName: "",
			Key:      "ct-q",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleDataUpdate,
		},
		{
			ViewName: "",
			Key:      "tab",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleViewLoop,
		},
		{
			ViewName: "note",
			Key:      "ct-s",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleSendNote,
		},
		{
			ViewName: "note-history",
			Key:      "up",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleHistorySelect("up"),
			// Handler:  controllers.CursorUp(),

		},
		{
			ViewName: "note-history",
			Key:      "down",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleHistorySelect("down"),
		},
		{
			ViewName: "note-history",
			Key:      "left",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleHistorySelect("left"),
		},
		{
			ViewName: "note-history",
			Key:      "right",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleHistorySelect("right"),
		},
		{
			ViewName: "note-history",
			Key:      "e",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleJumpToEnd(),
			// Handler:  controllers.CursorUp(),
		},
		{
			ViewName: "note-history",
			Key:      "delete",
			Modifier: gocui.ModNone,
			Handler:  gui.HandleDeleteNote,
		},
	}
}
