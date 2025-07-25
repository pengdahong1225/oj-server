package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"oj-server/module/configManager"
)

var (
	RedisClient *redis.Client
)

func Init() error {
	cfg := configManager.AppConf.RedisCfg
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
