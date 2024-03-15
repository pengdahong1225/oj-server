package internal

import (
	"judge-service/global"
	"sync"
)

type Server struct {
}

func (receiver Server) Start() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	err := global.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		consumerSrv := new(ConsumerServer)
		consumerSrv.startMQConsumer()
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}
