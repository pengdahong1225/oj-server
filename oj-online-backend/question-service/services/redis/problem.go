package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

// fields
const (
	UPStateField = "state"
)

// 用户提交状态缓存到redis，默认过期时间1分钟
func SetUPState(uid int64, problemID int64, state int) error {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)

	return SetKVByHashWithExpire(submitID, UPStateField, fmt.Sprintf("%d", state), 60)
}

func QueryUPState(uid int64, problemID int64) (int, error) {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)
	state, err := GetValueByHash(submitID, UPStateField)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(state)
}

func GetProblemHotData(problemID int64) (string, error) {
	conn := newConn()
	defer conn.Close()

	test, err := GetValueByHash(fmt.Sprintf("%d", problemID), "compile")
	if err != nil {
		return "", err
	}
	return test, nil
}

// QueryRankList 获取排行榜信息
func QueryRankList() ([]string, error) {
	conn := newConn()
	defer conn.Close()

	return redigo.Strings(conn.Do("zrange", "rank", 0, -1))
}

// SetJudgeResult 将结果写入redis，默认过期时间60分钟
func SetJudgeResult(uid int64, problemID int64, result string) error {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)
	return SetKVByStringWithExpire(submitID, result, 60*60)
}
func QueryJudgeResult(uid int64, problemID int64) (string, error) {
	return GetValueByString(fmt.Sprintf("%d:%d", uid, problemID))
}
