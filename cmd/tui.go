package cmd

import (
	"github.com/haochend413/mantis/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Init TUI",
	Long:  "Mantis TUI",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartTui()
	},
}
