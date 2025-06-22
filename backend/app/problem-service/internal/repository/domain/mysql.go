package db

import (
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"gorm.io/gorm"
)

var (
	DbSession *gorm.DB
	err       error
)

func Init() error {
	DbSession, err = db.NewMysqlSession()
	if err != nil {
		return err
	}

	return nil
}
