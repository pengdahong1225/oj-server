package redis

import (
	redigo "github.com/gomodule/redigo/redis"
)

// GetTestCaseJson 获取测试用例
func GetTestCaseJson(id int64) (string, error) {
	conn := newConn()
	defer conn.Close()

	ret, err := redigo.String(conn.Do("Get", id))
	if err != nil {
		return "", err
	}
	return ret, nil
}

// QueryRankList 获取排行榜信息
func QueryRankList() ([]string, error) {
	conn := newConn()
	defer conn.Close()

	return redigo.Strings(conn.Do("zrange", "rank", 0, -1))
}
