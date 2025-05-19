package ui

import (
	"github.com/spf13/cobra"
)

var UICmd = &cobra.Command{
	Use:   "ui",
	Short: "Gocui TUI",
	Long:  "Gocui TUI for Mantis",
	Run: func(cmd *cobra.Command, args []string) {
		RunUI()
	},
}
