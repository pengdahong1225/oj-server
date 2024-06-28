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

func GetValueByHash(key string, field string) (string, error) {
	conn := newConn()
	defer conn.Close()
	if reply, err := redis.String(conn.Do("HGet", key, field)); err != nil {
		return "", err
	} else {
		return reply, nil
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
