package cmd

import (
	"os"
	"strings"
)

// Execute with rootCmd.SetArgs(args)
func Execute(arg string) error {
	args := strings.Fields(arg)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	// //clean exit
	// if err == gocui.ErrQuit {
	// 	return gocui.ErrQuit
	// }
	if err != nil {
		os.Exit(1)
	}
	return nil
}
