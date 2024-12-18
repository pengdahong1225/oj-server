package daemon

import (
	"encoding/json"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/rpc/comment"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/svc/redis"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"github.com/pengdahong1225/oj-server/backend/module/mq"
	"github.com/pengdahong1225/oj-server/backend/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// Daemon 后台服务
type Daemon struct {
}

// CommentSaveConsumer
// comment消费者
func (receiver Daemon) CommentSaveConsumer() {
	consumer := mq.NewConsumer(
		consts.RabbitMqExchangeKind,
		consts.RabbitMqExchangeName,
		consts.RabbitMqCommentQueue,
		consts.RabbitMqCommentKey,
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	if deliveries == nil {
		logrus.Errorln("消费失败")
		return
	}
	defer consumer.Close()

	for d := range deliveries {
		if syncHandle(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false) // 拒绝并且丢失
		}
	}
}

func syncHandle(data []byte) bool {
	c := &pb.Comment{}
	err := proto.Unmarshal(data, c)
	if err != nil {
		logrus.Errorln("解析err：", err.Error())
		return false
	}
	// 处理
	goroutinePool.Instance().Submit(func() {
		comment.CommentSaveHandler{}.SaveComment(c)
	})

	return true
}

// RankList
// 排行榜维护
func (receiver Daemon) RankList() {
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
		logrus.Errorln("获取排行榜失败", result.Error.Error())
		return
	}

	if err := redis.UpdateRankList(orderList); err != nil {
		logrus.Errorln("更新排行榜失败", err.Error())
	}
}

// ReadTagList
// 提取tag list
func (receiver Daemon) ReadTagList() {
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
