package cmd

import (
	"fmt"

	"github.com/haochend413/mantis/app/state"
	"github.com/haochend413/mantis/defs"
	dailyui "github.com/haochend413/mantis/ui/dailyUI"
	"github.com/haochend413/mantis/ui/tui"

	"github.com/spf13/cobra"
)

var topicContent string
var ids []string

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Init noteUI",
	Long:  "Mantis TUI",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		appState := cmd.Context().Value(appStateKeyType{}).(*state.AppState)
		if cmd.Flags().Changed("add-topic") {
			// Flag is explicitly used
			//here we add task
			appState.DB_Data.TopicData = append(appState.DB_Data.TopicData, &defs.Topic{Topic: topicContent})
			if err := appState.DBManager.RefreshNoteTopic(appState.DB_Data); err != nil {
				fmt.Println("Failed to save topic to DB:", err)
			}
		} else if cmd.Flags().Changed("link") {
			//do the connection
			appState.DBManager.LinkNoteTopic(ids[0], ids[1])

		} else {
			// No flag provided, launch default UI
			tui.StartTui(appState)
		}

	},
}

var taskContent string

var dailyCmd = &cobra.Command{
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
		} else {
			// No flag provided, launch default UI
			dailyui.StartDailyUI(appState)
		}
	},
}

func init() {
	dailyCmd.Flags().StringVarP(&taskContent, "add", "a", "", "New Daily Task")
	noteCmd.Flags().StringVarP(&topicContent, "add-topic", "t", "", "New Topic")
	noteCmd.Flags().StringSliceVarP(&ids, "link", "l", nil, "Link Note with Topic")
}
