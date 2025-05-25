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
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
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

	//run cobra command
	if err := g.SetKeybinding("commandInput", gocui.KeyEnter, gocui.ModNone, sendCmd); err != nil {
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

// Keybindings for note-history
func noteHistoryKeys(g *gocui.Gui) {

	if err := g.SetKeybinding("noteHistory", gocui.KeyCtrlA, gocui.ModNone, quitNoteHistory); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("noteHistory", 'h', gocui.ModNone, CursorLeft); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("noteHistory", 'l', gocui.ModNone, CursorRight); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("noteHistory", 'j', gocui.ModNone, CursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("noteHistory", 'k', gocui.ModNone, CursorDown); err != nil {
		log.Panicln(err)
	}

}
