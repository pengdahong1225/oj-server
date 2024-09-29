package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"sync"
)

var (
	pool *redis.Pool
	once sync.Once
)

func init() {
	once.Do(func() {
		cfg := settings.Instance().RedisConfig
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
	})
}

func newConn() redis.Conn {
	return pool.Get()
}

func SetKVByString(key string, value string) error {
	conn := newConn()
	defer conn.Close()
	if _, err := conn.Do("Set", key, value); err != nil {
		return err
	}
	return nil
}

func SetKVByStringWithExpire(key string, value string, expire int) error {
	conn := newConn()
	defer conn.Close()
	if _, err := conn.Do("Set", key, value, "ex", expire); err != nil {
		return err
	}
	return nil
}

func GetValueByString(key string) (string, error) {
	conn := newConn()
	defer conn.Close()
	if reply, err := redis.String(conn.Do("Get", key)); err != nil {
		return "", err
	} else {
		return reply, nil
	}
}

func SetKVByHash(key string, field string, value string) error {
	conn := newConn()
	defer conn.Close()
	if _, err := conn.Do("HSet", key, field, value); err != nil {
		return err
	}
	return nil
}

// 要区分nil和error
func GetValueByHash(key string, field string) (string, error) {
	conn := newConn()
	defer conn.Close()

	reply, err := conn.Do("HGet", key, field)
	if err != nil {
		return "", err
	}
	if reply == nil {
		return "", nil // 不存在
	}

	if val, err := redis.String(reply, nil); err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func SetKVByHashWithExpire(key string, field string, value string, expire int) error {
	conn := newConn()
	defer conn.Close()
	if _, err := conn.Do("HSet", key, field, value, "ex", expire); err != nil {
		return err
	}
	return nil
}
