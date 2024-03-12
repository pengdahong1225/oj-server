package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"question-service/global"
)

type Amqp struct {
	MqConnection *amqp.Connection // 引用global的连接
	exchange     string
	queue        string
	routingKey   string
}

// CheckMqConnection 判断连接是否可用
func (receiver *Amqp) CheckMqConnection() {

}

func (receiver *Amqp) prepare(ctx *gin.Context) bool {
	conn := global.MqConnection
	// 判断连接是否可用

	// 建立通道
	ch, err := conn.Channel()
	if err != nil {
		logrus.Errorf("创建通道:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return false
	}
	// 声明交换机和队列
	err = ch.ExchangeDeclare(
		receiver.exchange, // 交换机名称
		"direct",          // 交换机类型
		true,              // 是否持久化
		false,             // 是否自动删除
		false,             // 是否独占
		false,             // 是否阻塞等待队列可用
		nil,               // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明交换机:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return false
	}
	queue, err := ch.QueueDeclare(
		receiver.queue, // 队列名称
		true,           // 是否持久化
		false,          // 是否自动删除
		false,          // 是否独占
		false,          // 是否阻塞等待队列可用
		nil,            // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("声明队列:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return false
	}

	err = ch.QueueBind(
		receiver.exchange,   // 交换机名称
		queue.Name,          // 队列名称
		receiver.routingKey, // 路由键
		false,               // 是否发送额外的bind headers
		nil,                 // 可选的额外参数
	)
	if err != nil {
		logrus.Errorf("绑定队列:%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "运行失败",
		})
		return false
	}
	return true
}

func (receiver *Amqp) publish() {

}
