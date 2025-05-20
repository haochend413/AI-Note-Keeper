package cmd

import (
	"github.com/haochend413/mantis/ui"
	"github.com/spf13/cobra"
)

// init command
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mantis",
	Short: "Mantis",
	Long:  "Maniis: Workflow manage tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ui.UIinit()
	},
}

// var UICmd = &cobra.Command{
// 	Use:   "ui",
// 	Short: "Gocui TUI",
// 	Long:  "Gocui TUI for Mantis",
// 	// Run: func(cmd *cobra.Command, args []string) {
// 	// 	RunUI()
// 	// },
// }
