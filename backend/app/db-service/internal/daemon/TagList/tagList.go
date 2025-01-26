package TagList

import (
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/sirupsen/logrus"
)

// ReadTagList
// 提取tag list
func ReadTagList() {
	db := mysql.Instance()
	var list [][]byte
	var tags []string

	/*
		SELECT tags FROM problem;
	*/
	result := db.Model(&mysql.Problem{}).Select("tags").Scan(&list)
	if result.Error != nil {
		logrus.Errorln("获取tag list失败", result.Error.Error())
		return
	}

	for _, val := range list {
		var temp []string
		if err := json.Unmarshal(val, &temp); err != nil {
			logrus.Errorln("解析tag list失败", err.Error())
			return
		}
		tags = append(tags, temp...)
	}
	if err := redis.UpdateTagList(tags); err != nil {
		logrus.Errorln("更新tag list失败", err.Error())
		return
	}
}
