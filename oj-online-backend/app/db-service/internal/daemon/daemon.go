package daemon

import (
	"encoding/json"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/rpc/comment"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/svc/mysql"
	"github.com/pengdahong1225/Oj-Online-Server/app/db-service/internal/svc/redis"
	"github.com/pengdahong1225/Oj-Online-Server/consts"
	"github.com/pengdahong1225/Oj-Online-Server/module/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/module/mq"
	"github.com/pengdahong1225/Oj-Online-Server/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// Daemon 后台服务
type Daemon struct {
}

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
	for d := range deliveries {
		if syncHandle(d.Body) {
			d.Ack(false)
		} else {
			d.Reject(false) // 拒绝并且丢失
		}
	}
}

// 解析，校验，提交任务给评测机
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

// 排行榜
func (receiver Daemon) LoopRank() {
	var key = "rank"
	conn := redis.NewConn()
	defer conn.Close()

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

	for _, item := range orderList {
		data, _ := json.Marshal(item.User) // 序列化
		_, err := conn.Do("ZADD", key, item.AccomplishCount, data)
		if err != nil {
			logrus.Errorln("更新排行榜失败", err.Error())
			break
		}
	}
}
