package cmd

import (
	"github.com/spf13/cobra"
)

// init command
// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "mantis",
	Short: "Mantis",
	Long:  "Mantis: Workflow manage tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// ui.UIinit()
		println("Welcome to Mantis!")
	},
}
