package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	AmqpClient *Client
	ExKind     string
	ExName     string
	QueName    string
	RoutingKey string
}

func (receiver *Producer) Publish(msg []byte) error {
	channel, err := receiver.AmqpClient.NewChannel(receiver.ExName, receiver.ExKind, receiver.QueName, receiver.RoutingKey)
	if err != nil {
		logrus.Errorf("获取channel失败, err: %s", err)
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		receiver.ExName,     // 指明交换机
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
		logrus.Errorf("发布消息失败:%s", err.Error())
		return err
	}
	return nil
}
