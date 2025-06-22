package cmd

import (
	"fmt"

	"github.com/haochend413/mantis/app/state"
	"github.com/haochend413/mantis/defs"
	dailyui "github.com/haochend413/mantis/ui/dailyUI"
	"github.com/haochend413/mantis/ui/tui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Init TUI",
	Long:  "Mantis TUI",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		appState := cmd.Context().Value(appStateKeyType{}).(*state.AppState)
		tui.StartTui(appState)
	},
}

var taskContent string

var dailyUICmd = &cobra.Command{
	Use:   "daily",
	Short: "Init DailyUI",
	Long:  "Mantis DailyUI",
	Args:  cobra.MaximumNArgs(1), // has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		appState := cmd.Context().Value(appStateKeyType{}).(*state.AppState)

		if cmd.Flags().Changed("add") {
			// Flag is explicitly used
			//here we add task
			appState.DB_Data.DailyTaskData = append(appState.DB_Data.DailyTaskData, &defs.DailyTask{Task: taskContent})
			if err := appState.DBManager.RefreshDaily(appState.DB_Data.DailyTaskData); err != nil {
				fmt.Println("Failed to save task to DB:", err)
			}
			dailyui.StartDailyUI(appState)
		} else {
			// No flag provided, launch default UI
			dailyui.StartDailyUI(appState)
		}
	},
}

func init() {
	dailyUICmd.Flags().StringVarP(&taskContent, "add", "a", "", "New Daily Task")
}
