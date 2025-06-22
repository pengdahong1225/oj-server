package db

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DbSession *gorm.DB
	Rdb       *redis.Client
	err       error
)

func Init() error {
	DbSession, err = db.NewMysqlSession()
	if err != nil {
		return err
	}

	Rdb = db.NewRedisClient()
	st := Rdb.Ping(context.Background())
	if st.Err() != nil {
		return st.Err()
	}

	return nil
}
