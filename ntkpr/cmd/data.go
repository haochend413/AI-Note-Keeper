package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/haochend413/ntkpr/config"
	"github.com/spf13/cobra"
)

var ExportNoteCmd = &cobra.Command{
	Use:   "export",
	Short: "export",
	Long:  "export",
	Run: func(cmd *cobra.Command, args []string) {
		if err := globalDB.ExportNoteToJSON(globalCfg.DataFilePath + "/notes.json"); err != nil {
			fmt.Fprintf(os.Stderr, "Export failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Exported notes to", globalCfg.DataFilePath+"/notes.json")
	},
}

var DataBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup ntkpr data",
	Long:  "Backup the ntkpr data folder to a specified destination, default to cwd.",
	Run: func(cmd *cobra.Command, args []string) {
		base, err := config.BasePathDefault()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting base path: %v\n", err)
			return
		}

		// timestamp folder name
		dest := fmt.Sprintf("ntkpr_backup_%s", time.Now().Format("2006-01-02_15-04-05"))
		if len(args) > 0 {
			dest = args[0]
		}

		// absolute
		if !filepath.IsAbs(dest) {
			cwd, _ := os.Getwd()
			dest = filepath.Join(cwd, dest)
		}

		cpCmd := exec.Command("cp", "-r", base, dest)
		if output, err := cpCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Error backing up: %v\n%s\n", err, output)
			return
		}

		fmt.Printf("Backed up %s to %s\n", base, dest)
	},
}
