package internal

import (
	"judge-service/internal/health"
	"judge-service/internal/logic"
	"judge-service/services/ants"
	"sync"
)

type Server struct {
}

func (receiver Server) Start() {
	wg := new(sync.WaitGroup)

	// 注册服务
	wg.Add(1)
	err := ants.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		healthSrv := new(health.HealthServer)
		healthSrv.Loop()
	})
	if err != nil {
		panic(err)
	}

	// 启动服务
	wg.Add(1)
	err = ants.AntsPoolInstance.Submit(func() {
		defer wg.Done()
		consumerSrv := new(logic.JudgeServer)
		consumerSrv.Loop()
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}
