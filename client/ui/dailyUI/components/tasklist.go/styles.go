package tasklist

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}
