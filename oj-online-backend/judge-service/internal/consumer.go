package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"judge-service/global"
)

type ConsumerServer struct {
}

// MQ消费者
func (receiver *ConsumerServer) startMQConsumer() {
	amqp := &Amqp{
		MqConnection: global.MqConnection,
		Exchange:     "amqp.direct",
		Queue:        "question",
		RoutingKey:   "question",
	}
	if err := amqp.Prepare(); err != nil {
		logrus.Errorf("amqp预处理失败:%s", err.Error())
		panic(err)
	}
	logrus.Infof("MQ消费者启动成功")
	defer amqp.Channel.Close()

	deliveries, err := amqp.Channel.Consume(
		amqp.Queue, // name
		"",         // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		logrus.Errorf("amqp错误:%s", err.Error())
		panic(err)
	}
	// 同步接收，异步处理
	for msg := range deliveries {
		// 异步处理
		global.AntsPoolInstance.Submit(func() {
			// 处理消息
			out := receiver.handleAsync(msg)
			// 回调
			receiver.callBack(out)
			// 处理完毕后，手动ack
			msg.Ack(true)
		})
	}
}

func (receiver *ConsumerServer) handleAsync(msg amqp.Delivery) []byte {
	return nil
}
func (receiver ConsumerServer) callBack(msg []byte) {

}
