/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "m",
	Short: "Mantis",
	Long:  "Mantis: Workflow manage tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hi")
	},
}

// Preboot setup
func init() {
	rootCmd.AddCommand(Quit)
}
