package dailydb

import (
	"github.com/haochend413/mantis/defs"
	"gorm.io/gorm"
)

// Sync DailyTaskData with the database
func (dd *DailyDB) SyncDailyTaskData(tasks []*defs.DailyTask) error {
	//This might be buggy: clear table
	dd.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&defs.DailyTask{})
	for _, n := range tasks {
		if result := dd.Db.Save(n); result.Error != nil {
			return result.Error
		}
	}
	return nil
}
