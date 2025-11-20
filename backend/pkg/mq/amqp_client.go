package mq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// amqp client
type Client struct {
	options *Options
	conn    *amqp.Connection
}
type Options struct {
	Host     string
	Port     int
	User     string
	PassWord string
	VHost    string
}

func NewClient(options *Options) *Client {
	return &Client{
		options: options,
	}
}

func (cli *Client) NewChannel(exName, exKind, quName, routingKey string) (*amqp.Channel, error) {
	if cli.conn == nil || cli.conn.IsClosed() {
		var err error
		cli.conn, err = cli.connect()
		if err != nil {
			return nil, err
		}
	}

	// 通道
	ch, err := cli.conn.Channel()
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	err = ch.QueueBind(
		queue.Name, // 队列名称
		routingKey, // 路由键
		exName,     // 交换机名称
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (cli *Client) connect() (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cli.options.User,
		cli.options.PassWord,
		cli.options.Host,
		cli.options.Port,
		cli.options.VHost,
	)
	conn, err := amqp.Dial(dsn)
	return conn, err
}
