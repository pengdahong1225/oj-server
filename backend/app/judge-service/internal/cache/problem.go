package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/registry"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

func UnLockUser(uid int64) error {
	key := fmt.Sprintf("lock:%d", uid)
	return rdb.Del(context.Background(), key).Err()
}

func SetJudgeResult(uid int64, problemID int64, result string, expire time.Duration) error {
	key := fmt.Sprintf("%d:%d:%s", uid, problemID, "result")
	return rdb.SetEx(context.Background(), key, result, expire).Err()
}

func SetUPState(uid int64, problemID int64, state int, expire time.Duration) error {
	key := fmt.Sprintf("%d:%d:%s", uid, problemID, "state")
	return rdb.SetEx(context.Background(), key, state, expire).Err()
}

func GetProblemConfig(problemID int64) (*pb.ProblemConfig, error) {
	problemConfig := &pb.ProblemConfig{}
	data, err := onGetProblemConfig(problemID)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(data, problemConfig); err != nil {
		return nil, err
	}
	return problemConfig, nil
}
func onGetProblemConfig(problemID int64) ([]byte, error) {
	data, err := rdb.HGet(context.Background(), strconv.FormatInt(problemID, 10), "hotData").Bytes()
	switch {
	case err == nil:
		return data, nil
	case errors.Is(err, redis.Nil):
		db, err := registry.NewDBConnection()
		if err != nil {
			logrus.Errorf("db服连接失败:%s\n", err.Error())
			return nil, err
		}
		defer db.Close()
		client := pb.NewProblemServiceClient(db)
		res, err := client.GetProblemHotData(context.Background(), &pb.GetProblemHotDataRequest{
			ProblemId: problemID,
		})
		if err != nil {
			logrus.Errorln(err.Error())
			return nil, err
		}
		return res.Data, nil
	}
	return nil, err
}
