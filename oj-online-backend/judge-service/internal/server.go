package internal

import (
	"github.com/sirupsen/logrus"
	"judge-service/services/goroutinePool"
	"judge-service/services/mq"
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
			handle(d.Body)
			d.Ack(true)
		})
	}
	wg.Wait()
}

// 解析，校验，提交任务给评测机
func handle(data []byte) {
}
