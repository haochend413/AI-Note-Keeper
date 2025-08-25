package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haochend413/mts/internal/app"
	"github.com/haochend413/mts/internal/db"
	"github.com/haochend413/mts/internal/ui"
)

func main() {
	// Initialize database
	dbConn, err := db.NewDB("notes.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	// Initialize application
	application := app.NewApp(dbConn)

	print(len(application.Topics))
	// Initialize UI model
	model := ui.NewModel(application)

	// Run Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
