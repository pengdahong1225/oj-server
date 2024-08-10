package mq

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Exname     string
	Exkind     string
	QuName     string
	RoutingKey string
	CTag       string
	channel    *amqp.Channel
}

func NewConsumer(exName, exKind, quName, routingKey, cTag string) *Consumer {
	return &Consumer{
		Exname:     exName,
		Exkind:     exKind,
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

func (receiver *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := receiver.channel.Cancel(receiver.CTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	logrus.Infoln("Consumer Channel shutdown OK")
	return nil
}
