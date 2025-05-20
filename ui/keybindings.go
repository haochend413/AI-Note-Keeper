package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

// Keybindings when back to init interface
func globalKeys(g *gocui.Gui) {
	// [ctrl-C] for exit ; Keybinding

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	// [:] for input bar setup;
	if err := g.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, setCommandInput); err != nil {
		log.Panicln(err)
	}

}

// Keybindings when focusing on cmdbars
func cmdinputKeys(g *gocui.Gui) {
	// [ctrl-C] for exit ; Keybinding

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("commandInput", gocui.KeyCtrlX, gocui.ModNone, quitCommandInput); err != nil {
		log.Panicln(err)
	}

}
