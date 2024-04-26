package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"judge-service/internal/logic"
	"judge-service/models"
	"judge-service/services/ants"
	"judge-service/services/mq"
	"judge-service/services/registry"
	"judge-service/settings"
	"net/http"
	"strings"
)

type JudgeConsumer struct {
}

func (receiver *JudgeConsumer) Loop() {
	amqp := &Amqp{
		MqConnection: mq.MqConnection,
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
		ants.AntsPoolInstance.Submit(func() {
			// 处理消息
			out := receiver.handleSync(msg)
			// 回调
			receiver.callBack(out)
			// 处理完毕后，手动ack
			msg.Ack(true)
		})
	}
}

func (receiver *JudgeConsumer) handleSync(msg amqp.Delivery) []byte {
	logrus.Infof("收到消息:%s", string(msg.Body))
	// 解析
	var form *models.JudgeRequest
	if err := json.Unmarshal(msg.Body, &form); err != nil {
		logrus.Errorf("解析消息失败:%s", err.Error())
		return nil
	}
	// 处理
	rspMsg := logic.Handle(form)
	return rspMsg
}

func (receiver *JudgeConsumer) callBack(msg []byte) {
	// 获取question服务地址
	dsn, err := registry.QuestionDsn(settings.Conf.RegistryConfig)
	if err != nil {
		logrus.Errorf("获取question服务地址失败:%s", err.Error())
		return
	}
	url := fmt.Sprintf("%s/%s", dsn, "questions/judgeCallback")
	body := strings.NewReader(string(msg))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logrus.Errorf("创建http请求失败:%s", err.Error())
		return
	}
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("发送http请求失败:%s", err.Error())
		return
	}
	defer rsp.Body.Close()
	logrus.Infof("回调结果:%s", rsp.Status)
}
