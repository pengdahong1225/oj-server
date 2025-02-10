package db

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	cfg := settings.Instance().RedisConfig
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	return rdb
}
