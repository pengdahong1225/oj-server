package cache

import (
	"context"
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func Init() error {
	cfg := settings.Instance().RedisConfig
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := RedisClient.Ping(context.Background())
	if st.Err() != nil {
		return st.Err()
	}
	return nil
}
