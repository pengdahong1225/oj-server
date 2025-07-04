package cache

import (
	"context"
	"fmt"
	"time"
)

const (
	LockUserKeyPrefix = "lock"
	LockUserTTL       = 60 * time.Second
)

func LockUser(uid int64) (bool, error) {
	key := fmt.Sprintf("%s:%d", LockUserKeyPrefix, uid)
	return rdb.SetNX(context.Background(), key, "locked", LockUserTTL).Result()
}
func UnLockUser(uid int64) error {
	key := fmt.Sprintf("%s:%d", LockUserKeyPrefix, uid)
	return rdb.Del(context.Background(), key).Err()
}
