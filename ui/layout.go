package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", maxX/4, maxY/4, maxX*3/4, maxY*3/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Main Window"
		fmt.Fprintln(v, "Welcome to Gocui!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func RunUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	//wait for close
	defer g.Close()

	//layout manager
	g.SetManagerFunc(layout)

	//ctrl-C for exit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}