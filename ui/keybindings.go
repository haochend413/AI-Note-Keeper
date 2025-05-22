package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

//This defines keybindings for different input states;

/* global keybindings */
func globalKeys(g *gocui.Gui) {
	// [ctrl-C] for exit ; Keybinding
	// g.DeleteKeybindings("commandInput")
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

// Keybindings for notekeys
func noteKeys(g *gocui.Gui) {
	// [ctrl-x] for input bar setup;
	// g.DeleteKeybindings("")

	//pull up command bar
	if err := g.SetKeybinding("note", gocui.KeyCtrlX, gocui.ModNone, setCommandInput); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("note", gocui.KeyCtrlA, gocui.ModNone, setNoteHistory); err != nil {
		log.Panicln(err)
	}

	//send note to noteDB, display on new view;
	if err := g.SetKeybinding("note", gocui.KeyEnter, gocui.ModNone, sendNote); err != nil {
		log.Panicln(err)
	}

}

// Keybindings when focusing on cmdbars
func cmdinputKeys(g *gocui.Gui) {
	// [ctrl-C] for exit ; Keybinding
	// g.DeleteKeybindings("")

	if err := g.SetKeybinding("commandInput", gocui.KeyCtrlX, gocui.ModNone, quitCommandInput); err != nil {
		log.Panicln(err)
	}
}

func noteHistoryKeys(g *gocui.Gui) {
	// [ctrl-C] for exit ; Keybinding
	// g.DeleteKeybindings("")

	if err := g.SetKeybinding("noteHistory", gocui.KeyCtrlA, gocui.ModNone, quitNoteHistory); err != nil {
		log.Panicln(err)
	}
}
