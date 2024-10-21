package cache

import (
	"context"
	"fmt"
	"time"
)

func LockUser(uid int64, expire int64) (bool, error) {
	key := fmt.Sprintf("lock:%d", uid)
	return Rdb.SetNX(context.Background(), key, "locked", time.Duration(expire)).Result()
}
func UnLockUser(uid int64) error {
	key := fmt.Sprintf("lock:%d", uid)
	return Rdb.Del(context.Background(), key).Err()
}
