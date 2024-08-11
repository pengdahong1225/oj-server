package mysql

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/common/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	db_  *gorm.DB
	once sync.Once
)

func Instance() *gorm.DB {
	once.Do(func() {
		cfg := settings.Instance().MysqlConfig
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User,
			cfg.Pwd, cfg.Host, cfg.Port, cfg.Db)
		var err error
		db_, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			// SkipDefaultTransaction: true, //全局禁用默认事务
		})
		if err != nil {
			panic(err)
		}
	})
	return db_
}
