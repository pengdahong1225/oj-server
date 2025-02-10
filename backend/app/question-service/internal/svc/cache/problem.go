package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"strconv"
)

// QueryRankList 获取排行榜信息
func QueryRankList() []string {
	return rdb.ZRange(context.Background(), "rank", 0, -1).Val()
}

func QueryUPState(uid int64, problemID int64) (int32, error) {
	key := fmt.Sprintf("%d:%d:%s", uid, problemID, "state")
	val, err := rdb.Get(context.Background(), key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		logrus.Infoln("key[%s]不存在", key)
		return -1, err
	case err != nil:
		logrus.Infoln("key[%s]查询错误:%s", key, err.Error())
		return -1, err
	case val == "":
		return -1, nil
	}
	state, _ := strconv.ParseInt(val, 10, 32)
	return int32(state), nil
}

func GetTagList() ([]string, error) {
	key := "tag_list"
	return rdb.SMembers(context.Background(), key).Result()
}

func GetJudgeResult(uid int64, problemID int64) (string, error) {
	key := fmt.Sprintf("%d:%d:%s", uid, problemID, "result")
	return rdb.Get(context.Background(), key).Result()
}
