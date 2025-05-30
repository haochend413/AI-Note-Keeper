package controllers

import "github.com/jroimartin/gocui"

func QuitApp() error {
	return gocui.ErrQuit
}
