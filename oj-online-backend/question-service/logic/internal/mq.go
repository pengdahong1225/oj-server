package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"question-service/global"
	"time"
)

type Amqp struct {
	MqConnection *amqp.Connection // 引用global的连接
	Exchange     string
	Queue        string
	RoutingKey   string
	Channel      *amqp.Channel // 通道
}

// CheckMqConnection 判断连接是否可用
func (receiver *Amqp) checkMqConnection() error {
	return nil
}

func (receiver *Amqp) Prepare() error {
	conn := global.MqConnection
	// 判断连接是否可用
	if err := receiver.checkMqConnection(); err != nil {
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

func (receiver *Amqp) Publish(msg []byte) bool {
	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/json",
		Body:         msg,
	}
	err := receiver.Channel.Publish(receiver.Exchange, receiver.RoutingKey, false, false, publishing)
	if err != nil {
		logrus.Errorf("发送消息失败:%s", err.Error())
		return false
	}
	logrus.Infof("发送消息成功:%s", msg)
	return true
}
