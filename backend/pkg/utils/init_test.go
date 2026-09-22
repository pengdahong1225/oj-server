package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"oj-server/global"
	"oj-server/pkg/db"
	"oj-server/proto/pb"
	"testing"
	"time"
)

func TestInitLeaderboard(t *testing.T) {
	// 初始化redis连接
	dsn := fmt.Sprintf("%s:%d", "192.168.2.111", 6379)
	rdb_, err := db.NewRedisCli(dsn)
	if err != nil {
		t.Errorf("redis init failed: %v", err)
		return
	}

	ctx := context.Background()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		userInfo := &pb.UserInfo{
			Uid:       int64(10000 + i),
			Nickname:  "test",
			AvatarUrl: "https://avatars.githubusercontent.com/u/106545107?v=4",
		}
		score := r.Intn(500)
		lb := &pb.LeaderboardUserInfo{
			Uid:      userInfo.Uid,
			UserName: userInfo.Nickname,
			Avatar:   userInfo.AvatarUrl,
			Score:    int32(score),
		}
		lbData, err := protojson.Marshal(lb)
		if err != nil {
			t.Errorf("protojson.Marshal failed: %v", err)
			return
		}
		rdb_.ZAdd(ctx, global.GetMonthLeaderboardKey(), redis.Z{Score: float64(lb.Score), Member: lbData})
		rdb_.ZAdd(ctx, global.GetDailyLeaderboardKey(), redis.Z{Score: float64(lb.Score), Member: lbData})
	}
	rdb_.Expire(ctx, global.GetMonthLeaderboardKey(), global.MonthLeaderboardTTL)
	rdb_.Expire(ctx, global.GetDailyLeaderboardKey(), global.DailyLeaderboardTTL)
}
