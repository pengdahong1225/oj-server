package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"strconv"
)

// fields
const (
	UPStateField  = "state"
	UPResultField = "result"
)

// 将用户提交状态写入redis，默认过期时间1分钟
func SetUPState(uid int64, problemID int64, state int) error {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)

	return SetKVByHashWithExpire(submitID, UPStateField, fmt.Sprintf("%d", state), 60)
}

func GetUPState(uid int64, problemID int64) (int, error) {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)
	state, err := GetValueByHash(submitID, UPStateField)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(state)
}

func SetUPResult(uid int64, problemID int64, result string) error {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)
	return SetKVByHashWithExpire(submitID, UPResultField, result, 60)
}

func GetUPResult(uid int64, problemID int64) (string, error) {
	submitID := fmt.Sprintf("%d:%d", uid, problemID)
	result, err := GetValueByHash(submitID, UPResultField)
	if err != nil {
		return "", err
	}
	return result, nil
}

func GetProblemCompileConfig(problemID int64) (string, error) {
	conn := newConn()
	defer conn.Close()

	test, err := GetValueByHash(fmt.Sprintf("%d", problemID), "compile")
	if err != nil {
		return "", err
	}
	return test, nil
}

func GetProblemRunConfig(problemID int64) (string, error) {
	conn := newConn()
	defer conn.Close()

	test, err := GetValueByHash(fmt.Sprintf("%d", problemID), "run")
	if err != nil {
		return "", err
	}
	return test, nil
}

func GetProblemTest(problemID int64) (string, error) {
	conn := newConn()
	defer conn.Close()
	return GetValueByHash(fmt.Sprintf("%d", problemID), "test")
}

// QueryRankList 获取排行榜信息
func QueryRankList() ([]string, error) {
	conn := newConn()
	defer conn.Close()

	return redigo.Strings(conn.Do("zrange", "rank", 0, -1))
}
