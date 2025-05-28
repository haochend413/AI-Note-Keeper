package models

import "github.com/jroimartin/gocui"

//This defines keybinder

type KeyBinder struct {
	ViewName string
	Key      string
	Modifier gocui.Modifier
	Handler  func() error
}
