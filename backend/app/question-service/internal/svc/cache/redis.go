package cache

import (
	"github.com/pengdahong1225/oj-server/backend/module/db"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = db.NewRedisClient()
}
