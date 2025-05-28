package models

import "github.com/jroimartin/gocui"

// Each Window defines one view
type Window struct {
	WindowName string
	View       *gocui.View
	OnDisplay  bool
}
