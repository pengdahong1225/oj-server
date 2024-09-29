package mq

import (
	"fmt"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var connection *amqp.Connection

func connect(cfg *settings.MqConfig) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.PassWord,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)
	logrus.Debugln(dsn)
	MqConnection, err := amqp.Dial(dsn)
	return MqConnection, err
}

func newChannel(exName, exKind, quName, routingKey string) *amqp.Channel {
	if connection == nil || connection.IsClosed() {
		var err error
		connection, err = connect(settings.Instance().MqConfig)
		if err != nil {
			logrus.Errorln(err)
			return nil
		}
	}

	// 通道
	ch, err := connection.Channel()
	if err != nil {
		logrus.Errorf("创建通道失败:%s", err.Error())
		return nil
	}
	err = ch.ExchangeDeclare(
		exName, // 交换机名称
		exKind, // 交换机类型
		true,   // 是否持久化
		false,  // 是否自动删除
		false,  // 是否独占
		false,  // 是否阻塞等待队列可用
		nil,    // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明交换机失败:%s", err.Error())
		return nil
	}
	queue, err := ch.QueueDeclare(
		quName, // 队列名称
		true,   // 是否持久化
		false,  // 是否自动删除
		false,  // 是否独占
		false,  // 是否阻塞等待队列可用
		nil,    // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明队列失败:%s", err.Error())
		return nil
	}
	err = ch.QueueBind(
		queue.Name, // 队列名称
		routingKey, // 路由键
		exName,     // 交换机名称
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		logrus.Errorf("绑定队列失败:%s", err.Error())
		return nil
	}

	return ch
}
