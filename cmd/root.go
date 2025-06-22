package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/haochend413/mantis/app/state"
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

type appStateKeyType struct{}

func Execute(appState *state.AppState) {
	//use context to pass app state
	ctx := context.Background()
	ctx = context.WithValue(ctx, appStateKeyType{}, appState)
	rootCmd.SetContext(ctx)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
