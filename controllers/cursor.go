package controllers

import (
	"github.com/haochend413/mantis/models"
	"github.com/jroimartin/gocui"
)

// This defines logic for cursor movements, make control functions;

func CursorMoveMaker(g *gocui.Gui, viewName string) func(d models.Direction) error {
	return func(d models.Direction) error {
		v, _ := g.View(viewName)
		switch d {
		case models.Up:
			return cursorUp(v)
		case models.Down:
			return cursorDown(v)
		case models.Left:
			return cursorLeft(v)
		case models.Right:
			return cursorRight(v)
		default:
			return nil
		}
	}
}

// This turns on / resets cursor
func CursorOn(g *gocui.Gui, view *gocui.View) error {
	g.Cursor = true
	lines := view.BufferLines()

	// // Remove trailing empty lines
	// for len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
	// 	lines = lines[:len(lines)-1]
	// }

	if len(lines) == 0 {
		return view.SetCursor(0, 0)
	}

	px := len(lines[len(lines)-1])
	py := len(lines) - 1
	return view.SetCursor(px, py)
}

// h-j-k-l defines cursor movements;
//right now: only for ineditable views; editable views need different modes that will be set later.

func cursorUp(view *gocui.View) error {
	//g.Cursor = true should have already been set
	//move up cursor

	//current position
	px, py := view.Cursor()
	if py != 0 {
		//get used to the way it handles errors!
		if err := view.SetCursor(px, py-1); err != nil {
			return err
		}
	}
	return nil

}

func cursorDown(view *gocui.View) error {
	//g.Cursor = true should have already been set
	//move up cursor

	//current position
	px, py := view.Cursor()
	//here, py should not be lower than the last line; -2 : trimmed empty line
	if py != len(view.BufferLines())-2 {
		//get used to the way it handles errors!
		if err := view.SetCursor(px, py+1); err != nil {
			return err
		}
	}
	return nil

}

func cursorLeft(view *gocui.View) error {
	//g.Cursor = true should have already been set
	//move up cursor

	//current position
	px, py := view.Cursor()
	if px != 0 {
		//get used to the way it handles errors!
		if err := view.SetCursor(px-1, py); err != nil {
			return err
		}
	}
	return nil

}

func cursorRight(view *gocui.View) error {
	//g.Cursor = true should have already been set
	//move up cursor

	//current position
	px, py := view.Cursor()
	line, err := view.Line(py)
	if err != nil || line == "" {
		return nil // either no such line, or it's empty
	}
	runes := []rune(line)
	if px < len(runes) {
		if err := view.SetCursor(px+1, py); err != nil {
			return err
		}
	}
	return nil

}
