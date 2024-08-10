package internal

import (
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/goroutinePool"
	"github.com/pengdahong1225/Oj-Online-Server/app/judge-service/services/mq"
	pb "github.com/pengdahong1225/Oj-Online-Server/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Server struct{}

func (receiver *Server) Loop() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln(err)
		}
	}()
	consumer := mq.NewConsumer(
		"judge",
		"direct",
		"judge",
		"judge",
		"", // 消费者标签，用于区别不同的消费者
	)
	deliveries := consumer.Consume()
	defer consumer.Shutdown()

	wg := new(sync.WaitGroup)
	for d := range deliveries {
		wg.Add(1)
		goroutinePool.PoolInstance.Submit(func() {
			defer wg.Done()
			if handle(d.Body) {
				d.Ack(false)
			} else {
				d.Reject(true)
			}
		})
	}
	wg.Wait()
}

// 解析，校验，提交任务给评测机
func handle(data []byte) bool {
	submitForm := &pb.SubmitForm{}
	err := proto.Unmarshal(data, submitForm)
	if err != nil {
		logrus.Errorln("解析消息遇到错误：", err.Error())
		return false
	}
	// ...

	return true
}
