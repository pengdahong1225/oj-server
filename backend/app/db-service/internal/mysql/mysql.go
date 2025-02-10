package mysql

import (
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"gorm.io/gorm"
)

var (
	DBSession *gorm.DB
	err       error
)

func Init() error {
	DBSession, err = db.NewMysqlSession()
	if err != nil {
		return err
	}
	return nil
}
