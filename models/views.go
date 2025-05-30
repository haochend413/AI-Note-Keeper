package models

import "github.com/jroimartin/gocui"

// Each Window defines one view

type Window struct {
	Name      string
	Title     string
	View      *gocui.View
	OnDisplay bool
	//define size
	X0 int
	Y0 int
	X1 int
	Y1 int
}
