package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	exKind     string
	exName     string
	queName    string
	routingKey string
	cTag       string
	channel    *amqp.Channel
}

func NewConsumer(exKind, exName, quName, routingKey, cTag string) *Consumer {
	return &Consumer{
		exKind:     exKind,
		exName:     exName,
		queName:    quName,
		routingKey: routingKey,
		cTag:       cTag,
		channel:    nil,
	}
}

func (receiver *Consumer) Consume() <-chan amqp.Delivery {
	channel := newChannel(receiver.exName, receiver.exKind, receiver.queName, receiver.routingKey)
	if channel == nil {
		logrus.Errorln("获取channel失败")
		return nil
	}
	receiver.channel = channel // 保存channel
	deliveries, err := receiver.channel.Consume(
		receiver.queName,
		receiver.cTag,
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

func (receiver *Consumer) Close() {
	receiver.channel.Close()
}
