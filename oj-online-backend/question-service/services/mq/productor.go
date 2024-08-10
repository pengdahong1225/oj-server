package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	Exname     string
	Exkind     string
	QuName     string
	RoutingKey string
}

func (receiver *Producer) Publish(msg []byte) bool {
	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/x-protobuf",
		Body:         msg,
	}

	channel := newChannel(receiver.Exname, receiver.Exkind, receiver.QuName, receiver.RoutingKey)
	if channel != nil {
		logrus.Errorln("获取channel失败")
		return false
	}
	defer channel.Close()

	err := channel.Publish(
		receiver.Exname,
		receiver.RoutingKey,
		false, // false 表示如果交换机无法找到符合条件的队列时消息会被丢弃
		false, // false 表示不需要立即被消费者接收
		publishing,
	)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return false
	}
	logrus.Infof("发送消息成功:%s", msg)
	return true
}
