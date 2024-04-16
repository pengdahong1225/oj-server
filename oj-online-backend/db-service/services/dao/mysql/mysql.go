package mysql

import (
	"db-service/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User,
		cfg.Pwd, cfg.Host, cfg.Port, cfg.Db)
	var e error
	DB, e = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})
	return e
}
