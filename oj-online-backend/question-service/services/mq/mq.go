package mq

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"question-service/settings"
)

var MqConnection *amqp.Connection // rabbitMQ全局连接

func Init(cfg *settings.MqConfig) error {
	var err error
	MqConnection, err = ConnectToMq(cfg)
	return err
}

func ConnectToMq(cfg *settings.MqConfig) (*amqp.Connection, error) {
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
