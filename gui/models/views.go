package models

import "github.com/jroimartin/gocui"

// Each Window defines one view
type Window struct {
	Name      string
	Title     string
	View      *gocui.View
	OnDisplay bool
}
