package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"judge-service/global"
)

type Amqp struct {
	MqConnection *amqp.Connection // 引用global的连接
	Exchange     string
	Queue        string
	RoutingKey   string
	Channel      *amqp.Channel // 通道
	Done         chan error
}

// CheckMqConnection 判断连接是否可用
func (receiver *Amqp) checkMqConnection() error {
	if receiver.MqConnection.IsClosed() {
		var err error
		// 重连
		global.MqConnection, err = global.ConnectToMq()
		if err != nil {
			return err
		}
	}
	return nil
}

func (receiver *Amqp) Prepare() error {
	conn := global.MqConnection
	// 判断连接是否可用
	if err := receiver.checkMqConnection(); err != nil {
		logrus.Errorf("MQ连接不可用:%s", err.Error())
		return err
	}
	// 建立通道
	ch, err := conn.Channel()
	if err != nil {
		logrus.Errorf("创建通道失败:%s", err.Error())
		return err
	}
	// 暂存通道
	receiver.Channel = ch
	// 声明交换机和队列
	err = ch.ExchangeDeclare(
		receiver.Exchange, // 交换机名称
		"direct",          // 交换机类型
		true,              // 是否持久化
		false,             // 是否自动删除
		false,             // 是否独占
		false,             // 是否阻塞等待队列可用
		nil,               // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明交换机失败:%s", err.Error())
		return err
	}
	queue, err := ch.QueueDeclare(
		receiver.Queue, // 队列名称
		true,           // 是否持久化
		false,          // 是否自动删除
		false,          // 是否独占
		false,          // 是否阻塞等待队列可用
		nil,            // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明队列失败:%s", err.Error())
		return err
	}

	err = ch.QueueBind(
		receiver.Exchange,   // 交换机名称
		queue.Name,          // 队列名称
		receiver.RoutingKey, // 路由键
		false,               // 是否发送额外的bind headers
		nil,                 // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("绑定队列失败:%s", err.Error())
		return err
	}
	return nil
}
