package cache

import (
	"context"
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Init() error {
	rdb = db.NewRedisClient()
	st := rdb.Ping(context.Background())
	if st.Err() != nil {
		return st.Err()
	}

	return nil
}
