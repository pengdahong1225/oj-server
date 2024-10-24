package logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type FileHook struct {
	errFile   *os.File
	infoFile  *os.File
	debugFile *os.File
	fileDate  string
	appName   string
	filePath  string
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

	err := os.MkdirAll(hook.filePath, os.ModePerm)
	if err != nil {
		return err
	}
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", hook.filePath, hook.appName, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", hook.filePath, hook.appName, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.%s.log.%s", hook.filePath, hook.appName, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

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
func InitLog(appName, path, level string) error {
	timer := time.Now().Format("2006-01-02")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// 日志文件
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.log.%s", path, "error", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.log.%s", path, "info", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.log.%s", path, "debug", timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	hook := &FileHook{
		errFile:   errFile,
		infoFile:  infoFile,
		debugFile: debugFile,
		fileDate:  timer,
		appName:   appName,
		filePath:  path,
	}
	logrus.AddHook(hook)
	if level == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	format := &nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		CallerFirst:     true,
	}
	logrus.SetFormatter(format)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	return nil
}
