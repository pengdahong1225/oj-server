package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"oj-server/module/configManager"
	"oj-server/module/db"
)

var (
	rdb *redis.Client
	err error
)

func Init() error {
	cfg := configManager.AppConf.RedisCfg
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb, err = db.NewRedisCli(dsn)
	if err != nil {
		return err
	}
	return nil
}
