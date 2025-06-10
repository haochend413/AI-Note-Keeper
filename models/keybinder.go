package models

import "github.com/awesome-gocui/gocui"

//This defines keybinder

type KeyBinder struct {
	ViewName string
	Key      string
	Modifier gocui.Modifier
	Handler  func(g *gocui.Gui, v *gocui.View) error
}
