package dbcontroller

import "github.com/haochend413/mantis/db"

type DBManager struct {
	DataBases *db.DataBases
}

func (m *DBManager) InitManager() error {
	m.DataBases = &db.DataBases{}
	m.DataBases.InitAll()
	return nil
}

func (m *DBManager) CloseManager() error {
	m.DataBases.CloseAll()
	return nil
}
