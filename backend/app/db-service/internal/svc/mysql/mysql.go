package mysql

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"sync"
	"time"
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

		path := fmt.Sprintf("%s/orm.log", settings.Instance().LogConfig.Path)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Errorf("gorm日志文件打开失败：%s", err.Error())
		}
		writer := io.MultiWriter(os.Stdout, file)
		newLogger := logger.New(
			log.New(writer, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,          // Don't include params in the SQL log
				Colorful:                  false,         // Disable color
			},
		)

		db_, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
			// SkipDefaultTransaction: true, //全局禁用默认事务
		})
		if err != nil {
			panic(err)
		}
	})
	return db_
}
