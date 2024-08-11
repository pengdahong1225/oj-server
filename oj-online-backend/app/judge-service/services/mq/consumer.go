package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Exkind     string
	Exname     string
	QuName     string
	RoutingKey string
	CTag       string
	channel    *amqp.Channel
}

func NewConsumer(exKind, exName, quName, routingKey, cTag string) *Consumer {
	return &Consumer{
		Exkind:     exKind,
		Exname:     exName,
		QuName:     quName,
		RoutingKey: routingKey,
		CTag:       cTag,
		channel:    nil,
	}
}

func (receiver *Consumer) Consume() <-chan amqp.Delivery {
	channel := newChannel(receiver.Exname, receiver.Exkind, receiver.QuName, receiver.RoutingKey)
	if channel != nil {
		logrus.Errorln("获取channel失败")
		return nil
	}
	receiver.channel = channel // 保存channel
	deliveries, err := receiver.channel.Consume(
		receiver.QuName,
		receiver.CTag,
		false, // 是否自动确认
		false, // 是否独占队列
		false, // true代表生产者和消费者不能是同一个connect
		false,
		nil,
	)
	if err != nil {
		logrus.Errorln("获取deliveries失败")
		return nil
	}
	return deliveries
}
