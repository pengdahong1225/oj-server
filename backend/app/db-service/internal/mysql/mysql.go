package mysql

import (
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"gorm.io/gorm"
)

var (
	DBSession *gorm.DB
	err       error
)

func init() {
	DBSession, err = db.NewMysqlSession()
	if err != nil {
		panic(err)
	}
}
