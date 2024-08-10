package daemon

import (
	"encoding/json"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/services/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/services/redis"
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

// Daemon 后台服务
type Daemon struct {
}

// 排行榜
func (receiver Daemon) loopRank() {
	var key = "rank"
	conn := redis.NewConn()
	defer conn.Close()

	db := mysql.DB
	var orderList []mysql.Statistics

	/**
	SELECT
	    user_info.*,
	    user_problem_statistics.accomplish_count
	FROM
	    user_problem_statistics
	JOIN
	    user_info ON user_problem_statistics.uid = user_info.id
	ORDER BY
	    user_problem_statistics.accomplish_count DESC
	LIMIT 50;
	*/
	result := db.Select("user_info.*, user_problem_statistics.accomplish_count").
		Joins("JOIN user_info ON user_problem_statistics.uid = user_info.id").
		Order("user_problem_statistics.accomplish_count desc").
		Limit(50).
		Find(&orderList)

	if result.Error != nil {
		logrus.Errorln("获取排行榜失败", result.Error.Error())
		return
	}

	for _, item := range orderList {
		data, _ := json.Marshal(item.User) // 序列化
		_, err := conn.Do("ZADD", key, item.AccomplishCount, data)
		if err != nil {
			logrus.Errorln("更新排行榜失败", err.Error())
			break
		}
	}
}
