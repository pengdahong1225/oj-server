package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consumer struct {
	AmqpClient *Client
	ExKind     string
	ExName     string
	QueName    string
	RoutingKey string
	CTag       string
	channel    *amqp.Channel
}

func (receiver *Consumer) Consume() <-chan amqp.Delivery {
	channel, err := receiver.AmqpClient.NewChannel(receiver.ExName, receiver.ExKind, receiver.QueName, receiver.RoutingKey)
	if err != nil {
		logrus.Errorf("获取channel失败, err: %s", err)
		return nil
	}
	receiver.channel = channel // 保存channel
	deliveries, err := receiver.channel.Consume(
		receiver.QueName,
		receiver.CTag,
		false, // 是否自动确认
		false, // 是否独占队列
		false, // true代表生产者和消费者不能是同一个connect
		false,
		nil,
	)
	if err != nil {
		logrus.Errorf("获取deliveries失败, err: %s", err)
		return nil
	}
	return deliveries
}

func (receiver *Consumer) Close() {
	receiver.channel.Close()
}
