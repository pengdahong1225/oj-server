package data

import (
	"context"
	"fmt"
	"io"
	"log"
	"oj-server/common/consts"
	"oj-server/svr/user/internal/config"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo struct {
	db_  *gorm.DB
	rdb_ *redis.Client
}

func NewRepo() (*Repo, error) {
	// mysql
	mysql_cfg := config.MyConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_cfg.User,
		mysql_cfg.Pwd, mysql_cfg.Host, mysql_cfg.Port, mysql_cfg.Db)
	timer := time.Now().Format("2006_01_02")
	filePath := fmt.Sprintf("%s/orm.%s.log", consts.LogPath, timer)
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

	// redis
	redis_cfg := config.MyConfig.Redis
	dsn = fmt.Sprintf("%s:%d", redis_cfg.Host, redis_cfg.Port)
	rdb_ := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb_.Ping(context.Background())
	if st.Err() != nil {
		return nil, st.Err()
	}

	return &Repo{
		db_:  db_,
		rdb_: rdb_,
	}, nil
}
