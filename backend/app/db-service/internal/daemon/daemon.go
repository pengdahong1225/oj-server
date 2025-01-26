package daemon

import (
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/daemon/TagList"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/daemon/comment"
	"github.com/pengdahong1225/oj-server/backend/app/db-service/internal/daemon/rankList"
	"github.com/pengdahong1225/oj-server/backend/module/goroutinePool"
	"sync"
)

// Daemon 后台服务
type Daemon struct {
}

func (receiver Daemon) Start() {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		comment.ConsumeComment()
	})

	wg.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		rankList.MaintainRankList()
	})

	wg.Add(1)
	goroutinePool.Instance().Submit(func() {
		defer wg.Done()
		TagList.ReadTagList()
	})

	wg.Wait()
}
