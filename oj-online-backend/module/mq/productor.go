package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	exName     string
	exKind     string
	queName    string
	routingKey string
}

func NewProducer(exKind, exName, quName, routingKey string) *Producer {
	return &Producer{
		exName:     exName,
		exKind:     exKind,
		queName:    quName,
		routingKey: routingKey,
	}
}

func (receiver *Producer) Publish(msg []byte) bool {
	channel := newChannel(receiver.exName, receiver.exKind, receiver.queName, receiver.routingKey)
	if channel == nil {
		logrus.Errorln("获取channel失败")
		return false
	}
	defer channel.Close()

	err := channel.Publish(
		receiver.exName,     // 指明交换机
		receiver.routingKey, // 指明路由键
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
