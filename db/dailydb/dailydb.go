package dailydb

import "gorm.io/gorm"

type DailyDB struct {
	Db *gorm.DB
	// Name string
}

// func (nd *NoteDB) Init() error {
// 	NoteDB.Db
// }

func (nd *DailyDB) Close() error {
	n, err := nd.Db.DB()
	if err != nil {
		return err
	}
	return n.Close()
}

func (nd *DailyDB) GetDB() *gorm.DB {
	return nd.Db
}
