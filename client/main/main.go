package main

import (
	"github.com/haochend413/mantis/app/state"
	"github.com/haochend413/mantis/cmd"
)

func main() {
	// init app state
	appState := state.NewAppState()
	cmd.Execute(appState)
}
