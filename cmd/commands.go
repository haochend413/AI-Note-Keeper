package cmd

import (
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

//As for my idea now, commands should also be specific to which view is currently open,
// which should be shown both in color (border highlight) and text (maybe bottom bar)

//global command

var Quit = &cobra.Command{
	Use:   "quit",
	Short: "Quit mantis",
	Long:  "Quit from Mantis",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gocui.ErrQuit
	},
}
