package logger

import (
	"testing"
)

func TestWorker_InitLog(t *testing.T) {
	appName := "test"
	path := "./"
	if err := Init(path, appName, DebugLevel); err != nil {
		t.Error(err.Error())
		return
	}
	Debugf("测试日志")
	Infof("测试日志")
	Errorf("测试日志")
}
