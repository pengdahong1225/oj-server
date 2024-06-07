package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"question-service/settings"
)

var pool *redis.Pool

func Init(cfg *settings.RedisConfig) error {
	dsn := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", dsn)
		},
		DialContext:     nil,
		TestOnBorrow:    nil,
		MaxIdle:         8,   // 最大空闲数
		MaxActive:       10,  // 最大连接数 -- 0表示无限制
		IdleTimeout:     100, // 最大空闲时间
		Wait:            false,
		MaxConnLifetime: 0,
	}
	return nil
}

func Close() {
	_ = pool.Close()
}

func newConn() redis.Conn {
	return pool.Get()
}

func SetKVByStringWithExpire(key string, value string, expire int) error {
	redisConn := newConn()
	defer redisConn.Close()
	if _, err := redisConn.Do("Set", key, value, "ex", expire); err != nil {
		return err
	}
	return nil
}
