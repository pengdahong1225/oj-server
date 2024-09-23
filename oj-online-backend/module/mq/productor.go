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
	channel := newChannel(receiver.Exname, receiver.Exkind, receiver.QuName, receiver.RoutingKey)
	err := channel.Publish(
		receiver.Exname,     // 指明交换机
		receiver.RoutingKey, // 指明路由键
		false,               // false 表示如果交换机无法找到符合条件的队列时消息会被丢弃
		false,               // false 表示不需要立即被消费者接收
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 消息持久化
			Timestamp:    time.Now(),
			ContentType:  "application/x-protobuf",
			Body:         msg,
		},
	)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return false
	}
	logrus.Infof("发送消息成功:%s", msg)
	return true
}
