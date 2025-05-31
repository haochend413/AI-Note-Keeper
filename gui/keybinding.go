package gui

import (
	"github.com/haochend413/mantis/gui/keybindings"
	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
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
			Handler:  HandleAppQuit,
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
	}
}
