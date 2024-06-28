package ants

import (
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"runtime"
)

var AntsPoolInstance *ants.Pool

func AntsPanicHandler(i interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logrus.Errorln("worker exits from panic: %s\n", string(buf[:n]))
}
