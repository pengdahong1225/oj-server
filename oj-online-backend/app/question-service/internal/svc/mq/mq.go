package mq

import (
	"fmt"
	"github.com/pengdahong1225/Oj-Online-Server/common/settings"
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

// 完整的声明交换机，声明队列，绑定
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
		true,   // 交换机持久化
		false,  // 是否自动删除(最后一个消费者取消订阅该队列时，自动删除队列)
		false,  // 是否独占(只能被一个消费者连接)
		false,  // 是否阻塞等待队列可用
		nil,    // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明交换机失败:%s", err.Error())
		return nil
	}
	queue, err := ch.QueueDeclare(
		quName, // 队列名称
		true,   // 队列持久化
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
		exName,     // 交换机名称
		queue.Name, // 队列名称
		routingKey, // 路由键
		false,      // 是否发送额外的bind headers
		nil,        // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("绑定队列失败:%s", err.Error())
		return nil
	}

	return ch
}