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

// 验证码
func SetCaptcha(mobile string, code string, expire int) error {
	key := fmt.Sprintf("captcha:%s", mobile)
	return Rdb.SetEx(context.Background(), key, code, time.Duration(expire)).Err()
}
func GetCaptcha(mobile string) (string, error) {
	key := fmt.Sprintf("captcha:%s", mobile)
	return Rdb.Get(context.Background(), key).Result()
}
