package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// FileHook 按天滚动写文件
type FileHook struct {
	mu       sync.Mutex
	logFile  *os.File
	fileDate string
	filePath string
	name     string
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	date := entry.Time.Format("2006_01_02")
	if hook.logFile == nil || hook.fileDate != date {
		if err := hook.rotate(date); err != nil {
			return err
		}
	}
	return hook.write(entry)
}

func (hook *FileHook) rotate(date string) error {
	if err := os.MkdirAll(hook.filePath, 0755); err != nil {
		return fmt.Errorf("create log folder %s failed: %v", hook.filePath, err)
	}

	logFile, err := os.OpenFile(
		fmt.Sprintf("%s/%s.%s.log", hook.filePath, hook.name, date),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644,
	)
	if err != nil {
		return fmt.Errorf("open log file of %s failed: %v", date, err)
	}

	if hook.logFile != nil {
		_ = hook.logFile.Close()
	}
	hook.logFile = logFile
	hook.fileDate = date
	return nil
}

func (hook *FileHook) write(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	_, err = hook.logFile.Write(line)
	return err
}

// Close 关闭当前日志文件，供进程优雅退出时调用。
func (hook *FileHook) Close() error {
	hook.mu.Lock()
	defer hook.mu.Unlock()
	if hook.logFile != nil {
		err := hook.logFile.Close()
		hook.logFile = nil
		return err
	}
	return nil
}
