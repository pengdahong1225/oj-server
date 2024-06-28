package redis

import (
	"question-service/models"
	"strconv"
)

// fields
const (
	UserStateField = "state"
)

func GetUserState(uid int64) (int, error) {
	state, err := GetValueByHash(strconv.FormatInt(uid, 10), UserStateField)
	if err != nil {
		return -1, err
	}
	if state == "" {
		return models.UserStateNormal, nil
	}
	return strconv.Atoi(state)
}

func SetUserState(uid int64, state int) error {
	return SetKVByHash(strconv.FormatInt(uid, 10), UserStateField, strconv.Itoa(state))
}
