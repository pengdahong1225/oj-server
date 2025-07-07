package cache

import (
	"context"
	"fmt"
	"time"
)

func SetTaskState(taskId string, state int, expire time.Duration) error {
	key := fmt.Sprintf("%s:%s", taskId, "state")
	return rdb.SetEx(context.Background(), key, state, expire).Err()
}

func UnLockUser(uid int64) error {
	key := fmt.Sprintf("lock:%d", uid)
	return rdb.Del(context.Background(), key).Err()
}

func SetJudgeResult(taskId, result string, expire time.Duration) error {
	key := fmt.Sprintf("%s_%s", taskId, "result")
	return rdb.SetEx(context.Background(), key, result, expire).Err()
}
