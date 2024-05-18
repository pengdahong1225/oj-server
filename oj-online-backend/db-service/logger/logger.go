package logger

import (
	"db-service/settings"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
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

	app, err := settings.GetSystemConf("db-service")
	if err != nil {
		return err
	}
	path := settings.Conf.LogConfig.Path
	name := app.Name
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)

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
func Init() error {
	// 目录
	app, err := settings.GetSystemConf("db-service")
	if err != nil {
		return err
	}
	path := settings.Conf.LogConfig.Path
	name := app.Name
	timer := time.Now().Format("2006-01-02")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// 日志文件
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", path, name, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)

	hook := &FileHook{
		errFile:   errFile,
		infoFile:  infoFile,
		debugFile: debugFile,
		fileDate:  timer,
	}
	logrus.AddHook(hook)
	if settings.Conf.LogConfig.Level == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	return nil
}
