package cache

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	Rdb  *redis.Client
	once sync.Once
)

func init() {
	once.Do(func() {
		cfg := settings.Instance().RedisConfig
		dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		Rdb = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    dsn,
		})
	})
}
