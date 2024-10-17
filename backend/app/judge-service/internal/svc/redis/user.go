package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func LockUser(uid int64, expire int64) (bool, error) {
	conn := newConn()
	defer conn.Close()

	key := fmt.Sprintf("lock:%d", uid)
	return redis.Bool(conn.Do("SetNx", key, "locked", "ex", expire))
}
func UnLockUser(uid int64) (bool, error) {
	conn := newConn()
	defer conn.Close()

	key := fmt.Sprintf("lock:%d", uid)
	return redis.Bool(conn.Do("DEL", key))
}
