package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type FileHook struct {
	// mu serializes Fire. logrus fires hooks without holding its own lock, so
	// the hook has to guard its state by itself.
	mu       sync.Mutex
	logFile  *os.File
	fileDate string
	filePath string
	name     string
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire writes one entry to the log file of the entry's own day, rotating the
// file when the day changes.
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

// rotate switches the hook to the log file of date. It opens the new file
// before closing the previous one, so that a failed rotation keeps the current
// handle usable instead of leaving the hook without a file.
func (hook *FileHook) rotate(date string) error {
	if err := os.MkdirAll(hook.filePath, os.ModePerm); err != nil {
		return fmt.Errorf("create log folder %s failed: %v", hook.filePath, err)
	}

	logFile, err := os.OpenFile(
		fmt.Sprintf("%s/%s.%s.log", hook.filePath, hook.name, date),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666,
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

	if _, err = hook.logFile.Write(line); err != nil {
		return err
	}
	return nil
}

// Init initializes logger
func Init(filePath, name string, level logrus.Level) error {
	timer := time.Now().Format("2006_01_02")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	// log file
	// A failure is not fatal here: the hook retries on the next Fire call.
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.%s.log", filePath, name, timer), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open log file failed: %v\n", err)
	}
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
