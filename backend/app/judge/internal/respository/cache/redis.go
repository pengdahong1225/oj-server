package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"oj-server/module/configManager"
)

var (
	rdb *redis.Client
)

func Init() error {
	cfg := configManager.Instance().RedisCfg
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb = redis.NewClient(&redis.Options{
		Addr:    dsn,
		Network: "tcp",
	})
	st := rdb.Ping(context.Background())
	if st.Err() != nil {
		return st.Err()
	}
	return nil
}
