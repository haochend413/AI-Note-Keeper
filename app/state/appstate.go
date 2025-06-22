package state

import (
	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/defs"
)

// Data entities that are shared across the app
type AppState struct {
	//Cached db data
	DB_Data defs.DB_Data
	//interact with actual db using cached data
	DBManager *dbcontroller.DBManager
}

func NewAppState() *AppState {
	manager := &dbcontroller.DBManager{}
	manager.InitManager()
	DB_Data := manager.FetchAll()
	return &AppState{
		DB_Data:   *DB_Data,
		DBManager: manager,
	}
}
