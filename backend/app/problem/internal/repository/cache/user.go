package cache

import (
	"context"
	"fmt"
	"time"
)

func LockUser(uid int64, expire time.Duration) (bool, error) {
	key := fmt.Sprintf("lock:%d", uid)
	return rdb.SetNX(context.Background(), key, "locked", expire).Result()
}
func UnLockUser(uid int64) error {
	key := fmt.Sprintf("lock:%d", uid)
	return rdb.Del(context.Background(), key).Err()
}

// SetCaptcha 验证码
func SetCaptcha(mobile string, code string, expire int) error {
	key := fmt.Sprintf("captcha:%s", mobile)
	return rdb.SetEx(context.Background(), key, code, time.Duration(expire)).Err()
}
func GetCaptcha(mobile string) (string, error) {
	key := fmt.Sprintf("captcha:%s", mobile)
	return rdb.Get(context.Background(), key).Result()
}
