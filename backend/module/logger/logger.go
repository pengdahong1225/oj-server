package logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

type FileHook struct {
	logFile  *os.File
	fileDate string
	filePath string
	name     string
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
	hook.logFile.Close()

	err := os.MkdirAll(hook.filePath, os.ModePerm)
	if err != nil {
		return err
	}
	logFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.log.%s", hook.filePath, hook.name, timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	hook.logFile = logFile
	hook.fileDate = timer

	hook.write(entry)

	return nil
}

func (hook FileHook) write(entry *logrus.Entry) {
	line, _ := entry.Bytes()
	_, _ = hook.logFile.Write(line)
}

// Init 初始化日志
func Init(filePath, name string, level logrus.Level) error {
	timer := time.Now().Format("2006-01-02")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	// 日志文件
	logFile, _ := os.OpenFile(fmt.Sprintf("%s/%s.log.%s", filePath, name, timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	hook := &FileHook{
		logFile:  logFile,
		fileDate: timer,
		filePath: filePath,
		name:     name,
	}
	logrus.AddHook(hook)
	logrus.SetLevel(level)

	format := &nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			filename := path.Base(f.File)
			return fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	}

	logrus.SetFormatter(format)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	return nil
}
