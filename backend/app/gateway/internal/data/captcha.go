package data

import (
	"context"
	"fmt"
	"oj-server/global"
)

func SetSmsCaptcha(phone string, value string) error {
	key := fmt.Sprintf("%s:%s", global.CaptchaPrefix, phone)
	return rdb.SetEx(context.Background(), key, value, global.CaptchaExpired).Err()
}
func GetSmsCaptcha(phone string) (string, error) {
	key := fmt.Sprintf("%s:%s", global.CaptchaPrefix, phone)
	return rdb.Get(context.Background(), key).Result()
}
