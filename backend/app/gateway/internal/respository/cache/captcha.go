package cache

import (
	"context"
	"fmt"
	"time"
)

// ============================短信验证码============================

const (
	prefix = "captcha"
	expire = 60 * time.Second
)

func SetImageCaptcha(phone string, value string) error {
	key := fmt.Sprintf("%s:%s", prefix, phone)
	return rdb.SetEx(context.Background(), key, value, expire).Err()
}
func GetImageCaptcha(phone string) (string, error) {
	key := fmt.Sprintf("%s:%s", prefix, phone)
	return rdb.Get(context.Background(), key).Result()
}
