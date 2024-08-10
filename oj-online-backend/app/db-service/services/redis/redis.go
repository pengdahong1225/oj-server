package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pengdahong1225/Oj-Online-Server/config"
)

var pool *redis.Pool

func Init(cfg *config.RedisConfig) error {
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

func NewConn() redis.Conn {
	return pool.Get()
}

func SetKVByHash(key string, field string, value string) error {
	conn := NewConn()
	defer conn.Close()
	if _, err := conn.Do("HSet", key, field, value); err != nil {
		return err
	}
	return nil
}
