package keybindings

import (
	"github.com/haochend413/mantis/gui/models"
	"github.com/jroimartin/gocui"
)

// Maybe just do ordinary handlers, and then do abstractions later;
// Maybe a better way is to define multiple keybindings in the global view, including opening up and shutting down existing views.

func CreateAllKeybinders() []*models.KeyBinder {
	return []*models.KeyBinder{
		{
			ViewName: "",
			Key:      "ct-c",
			Modifier: gocui.ModNone,
			Handler:  HandleQuit,
		},
	}
}
