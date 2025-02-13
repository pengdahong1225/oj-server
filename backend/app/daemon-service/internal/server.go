package internal

import (
	ServerBase "github.com/pengdahong1225/oj-server/backend/app/common/serverBase"
	"github.com/pengdahong1225/oj-server/backend/app/daemon-service/internal/rank_list"
	"github.com/pengdahong1225/oj-server/backend/app/daemon-service/internal/tag_list"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"sync"
)

type Server struct {
	ServerBase.Server
}

func (receiver *Server) Start() {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	_ = goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		rank_list.MaintainRankList()
	})

	wg.Add(1)
	_ = goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		tag_list.ReadTagList()
	})

	// 注册服务节点
	err := receiver.Register()
	if err != nil {
		panic(err)
	}
	wg.Wait()
}
