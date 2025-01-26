package rankList

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/sirupsen/logrus"
	"time"
)

// MaintainRankList
// 排行榜维护
func MaintainRankList() {
	timer := time.NewTimer(getDuration())
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			refresh()
			timer.Reset(getDuration())
		}
	}
}
func getDuration() time.Duration {
	now := time.Now()
	// 计算明天零点的时间
	tomorrow := now.Add(24 * time.Hour)
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())

	duration := midnight.Sub(now)
	if duration < 0 {
		// 如果已经过了今天的零点，则调整为第二天零点
		duration += 24 * time.Hour
	}

	return duration
}
func refresh() {
	db := mysql.Instance()
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
		logrus.Errorf("获取排行榜失败, err=%s", result.Error.Error())
		return
	}

	if err := redis.UpdateRankList(orderList); err != nil {
		logrus.Errorf("更新排行榜失败, err=%s", err.Error())
	}
}
