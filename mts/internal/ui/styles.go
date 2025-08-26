package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/table"
)

// UI styles
var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("123"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("211")).
			Bold(true).
			Padding(0, 1)

	// topicStyle = lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("86")).
	// 		Background(lipgloss.Color("235")).
	// 		Padding(0, 1).
	// 		Margin(0, 1, 0, 0)

	simpleTopicsStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.HiddenBorder()).
				Padding(0, 0).
				Margin(0, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0, 0, 2)
)

var baseTableStyle = table.Styles{
	Header: lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color("252")),

	Cell: lipgloss.NewStyle().
		Padding(0, 1),
	// Foreground(lipgloss.Color("246")),
	Selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("49")), // yellow
	// 	Background(lipgloss.Color("236")), // dark gray

}

var focusedTableStyle = table.Styles{
	Header: lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1),
	// Foreground(lipgloss.Color("229")),

	Cell: lipgloss.NewStyle().
		Padding(0, 1),
	Selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("49")).
		Background(lipgloss.Color("243")). // purple
		Bold(true),
}
