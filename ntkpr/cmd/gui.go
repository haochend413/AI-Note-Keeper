package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var LaunchGUICmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch Gui. ",
	Long:  "Launch the Gui for lastest exported notes.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if node is installed
		nodeCheck := exec.Command("node", "--version")
		if output, err := nodeCheck.CombinedOutput(); err != nil {
			fmt.Println("Node.js is not installed. Attempting to install...")

			// Try to install using brew (macOS) or common Linux package managers
			var installNode *exec.Cmd
			if _, err := exec.LookPath("brew"); err == nil {
				installNode = exec.Command("brew", "install", "node")
			} else if _, err := exec.LookPath("apt"); err == nil {
				installNode = exec.Command("sudo", "apt", "install", "-y", "nodejs", "npm")
			} else if _, err := exec.LookPath("yum"); err == nil {
				installNode = exec.Command("sudo", "yum", "install", "-y", "nodejs", "npm")
			} else {
				fmt.Fprintf(os.Stderr, "Could not find a package manager to install Node.js\n")
				fmt.Fprintf(os.Stderr, "Please install manually from https://nodejs.org/\n")
				return
			}

			installNode.Stdout = os.Stdout
			installNode.Stderr = os.Stderr
			if err := installNode.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to install Node.js: %v\n", err)
				return
			}
			fmt.Println("Node.js installed successfully!")
		} else {
			fmt.Printf("Node.js version: %s", string(output))
		}

		// Check if pnpm is installed
		pnpmCheck := exec.Command("pnpm", "--version")
		if output, err := pnpmCheck.CombinedOutput(); err != nil {
			fmt.Println("pnpm is not installed. Attempting to install...")

			installPnpm := exec.Command("npm", "install", "-g", "pnpm")
			installPnpm.Stdout = os.Stdout
			installPnpm.Stderr = os.Stderr
			if err := installPnpm.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to install pnpm: %v\n", err)
				return
			}
			fmt.Println("pnpm installed successfully!")
		} else {
			fmt.Printf("pnpm version: %s", string(output))
		}

		guiDir, err := filepath.Abs("../gui")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving gui path: %v\n", err)
			return
		}

		cmd1 := exec.Command("pnpm", "install")
		cmd1.Dir = guiDir
		cmd1.Stdout = os.Stdout
		cmd1.Stderr = os.Stderr
		if err := cmd1.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing pnpm packages: %v\n", err)
			return
		}

		cmd2 := exec.Command("pnpm", "dev")
		cmd2.Dir = guiDir
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		fmt.Println("Starting GUI at:", guiDir)
		if err := cmd2.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running pnpm dev: %v\n", err)
			return
		}
	},
}
