package go_pool

import (
	"runtime"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

var (
	pool *ants.Pool
	once sync.Once
)

func init() {
	once.Do(func() {
		var err error
		pool, err = ants.NewPool(1024*10, ants.WithPanicHandler(panicHandler))
		if err != nil {
			panic(err)
		}
	})
}

func panicHandler(i interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logrus.Errorf("worker exits from panic: %s\n", string(buf[:n]))
}

func Submit(task func()) error {
	return pool.Submit(task)
}
