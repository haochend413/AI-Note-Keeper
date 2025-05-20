package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// this defines the major layout
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", 0, 0, maxX/4, maxY-10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Main Window"
		fmt.Fprintln(v, "Welcome to Gocui!")
	}
	if u, err := g.SetView("second", maxX*2/4, 0, maxX, maxY-10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		u.Title = "Second Main Window"
		fmt.Fprintln(u, "Welcome to Gocui!")
	}
	return nil
}

// this function inits the major UI interface when rootcmd is called.
func UIinit() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	//check startup err
	if err != nil {
		log.Panicln(err)
	}
	//wait for close
	defer g.Close()

	//layout manager
	g.SetManagerFunc(layout)

	//setup keybindings;
	globalKeys(g)
	// //ctrl-C for exit ; Keybinding
	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
