package internal

import (
	"judge-service/global"
	"sync"
)

type Server struct {
}

func (receiver Server) Start() {
	wg := new(sync.WaitGroup)

	// 注册服务
	wg.Add(1)
	err := global.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		healthSrv := new(HealthServer)
		healthSrv.start()
	})
	if err != nil {
		panic(err)
	}

	// 启动服务
	wg.Add(1)
	err = global.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		consumerSrv := new(ConsumerServer)
		consumerSrv.start()
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}
