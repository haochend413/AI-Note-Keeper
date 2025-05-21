/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/haochend413/mantis/cmd"
	"github.com/haochend413/mantis/db"
	"github.com/haochend413/mantis/ui"
)

func main() {
	//init note database
	db.NoteDBInit()
	//init Cobra
	cmd.Execute()
	//init Gocui
	ui.UIinit()
}
