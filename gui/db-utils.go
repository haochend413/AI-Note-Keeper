package gui

import (
	"github.com/haochend413/mantis/models"
)

// refresh the layout when db has been changed
// func (gui *Gui) FetchNewestData() error {

// this stored all the temporary db data that will be used to store & update DB when the app init / close or when a specific function requires database operations;
var DB_Data *models.DB_Data = &models.DB_Data{}
