package models

import "github.com/jroimartin/gocui"

// Each Window defines one view

type Window struct {
	Name  string
	Title string
	View  *gocui.View
	//define window size
	X0 int
	Y0 int
	X1 int
	Y1 int
	//define view configs
	OnDisplay bool
	Editable  bool
	Scroll    bool
	Cursor    bool
	// Wrap     bool
}
