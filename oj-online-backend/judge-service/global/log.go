package global

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type FileHook struct {
	errFile   *os.File
	infoFile  *os.File
	debugFile *os.File
	fileDate  string
}

func (hook FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 按时间和级别写入
func (hook FileHook) Fire(entry *logrus.Entry) error {
	timer := entry.Time.Format("2006-01-02")
	if hook.fileDate == timer {
		hook.write(entry)
		return nil
	}
	hook.errFile.Close()
	hook.infoFile.Close()
	hook.debugFile.Close()

	path := ConfigInstance.Log_.Path
	app := ConfigInstance.System_.Name
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)

	hook.errFile = errFile
	hook.infoFile = infoFile
	hook.debugFile = debugFile
	hook.fileDate = timer
	hook.write(entry)
	return nil
}

func (hook FileHook) write(entry *logrus.Entry) {
	line, _ := entry.String()
	switch entry.Level {
	case logrus.ErrorLevel:
		hook.errFile.Write([]byte(line))
	case logrus.InfoLevel:
		hook.infoFile.Write([]byte(line))
	case logrus.DebugLevel:
		hook.debugFile.Write([]byte(line))
	}
}

// 初始化日志
func initLog() error {
	// 目录
	path := ConfigInstance.Log_.Path
	app := ConfigInstance.System_.Name
	timer := time.Now().Format("2006-01-02")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// 日志文件
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, app, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)

	hook := &FileHook{
		errFile:   errFile,
		infoFile:  infoFile,
		debugFile: debugFile,
		fileDate:  timer,
	}
	logrus.AddHook(hook)
	if ConfigInstance.Log_.Level == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	return nil
}
