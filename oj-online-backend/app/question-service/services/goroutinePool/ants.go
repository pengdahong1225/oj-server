package goroutinePool

import (
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"runtime"
)

var PoolInstance *ants.Pool

func PanicHandler(i interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logrus.Errorln("worker exits from panic: %s\n", string(buf[:n]))
}

func Init() error {
	var err error
	PoolInstance, err = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(PanicHandler))
	return err
}
