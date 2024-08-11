package goroutinePool

import (
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
)

var (
	pool *ants.Pool
	once sync.Once
)

func Instance() *ants.Pool {
	once.Do(func() {
		var err error
		pool, err = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(panicHandler))
		if err != nil {
			panic(err)
		}
	})
	return pool
}

func panicHandler(i interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logrus.Errorln("worker exits from panic: %s\n", string(buf[:n]))
}
