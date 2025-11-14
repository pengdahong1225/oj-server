package logger

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestWorker_InitLog(t *testing.T) {
	appName := "gateway"
	path := "./"
	if err := Init(path, appName, logrus.InfoLevel); err != nil {
		t.Error(err.Error())
		return
	}
	logrus.Debugf("测试日志")
	logrus.Infof("测试日志")
	logrus.Errorf("测试日志")
}
