package handler

import (
	"db-service/global"
	"db-service/internal/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

func StartDaemon() {
	// 周期性定时器
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			Daemon{}.loopRank()
		}
	}
}

type Daemon struct {
}

// 维护排行榜
func (receiver Daemon) loopRank() {
	var key = "rank"
	conn := global.RedisPoolInstance.Get()
	defer conn.Close()

	// select phone,nickname,pass_count from user_info order by pass_count desc;
	var users []models.UserInfo
	result := global.DBInstance.Select("phone", "nickname", "pass_count").Order("pass_count desc").Limit(50).Find(&users)
	if result.Error != nil {
		logrus.Errorln("获取排行榜失败", result.Error.Error())
		return
	}

	for _, user := range users {
		data, _ := json.Marshal(user) // 序列化
		_, err := conn.Do("ZADD", key, user.PassCount, data)
		if err != nil {
			logrus.Errorln("更新排行榜失败", err.Error())
			break
		}
	}
}
