package redis

import (
	redigo "github.com/gomodule/redigo/redis"
)

// 获取测试用例
func GetTestCaseJson(id int64) (string, error) {
	conn := NewConn()
	defer conn.Close()

	ret, err := redigo.String(conn.Do("Get", id))
	if err != nil {
		return "", err
	}
	return ret, nil
}
