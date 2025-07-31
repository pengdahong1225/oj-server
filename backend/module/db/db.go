package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"
)

func NewMysqlCli(dsn, logPath string) (*gorm.DB, error) {
	filePath := fmt.Sprintf("%s/orm.log", logPath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
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

	db_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		// SkipDefaultTransaction: true, //全局禁用默认事务
	})
	if err != nil {
		return nil, err
	}

	return db_, nil
}

func NewRedisCli(dsn string) (*redis.Client, error) {
	rdb_ := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb_.Ping(context.Background())
	if st.Err() != nil {
		return nil, st.Err()
	}
	return rdb_, nil
}
