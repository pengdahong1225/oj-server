package redis

import (
	"context"
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func UpdateRankList(orderList []mysql.Statistics) error {
	var key = "rank"
	for _, item := range orderList {
		data, _ := json.Marshal(item.User) // 序列化
		err := rdb.ZAdd(context.Background(), key, redis.Z{
			Score:  float64(item.AccomplishCount),
			Member: data,
		}).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// CacheProblemConfig 缓存题目热点数据
func CacheProblemConfig(problemID int64, data []byte) error {
	return rdb.HSet(context.Background(), strconv.FormatInt(problemID, 10), "hotData", data).Err()
}