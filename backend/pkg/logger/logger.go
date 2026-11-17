package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	std  *logrus.Logger
	once sync.Once
)

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
)

func Init(filePath, name string, level logrus.Level) error {
	var initErr error
	once.Do(func() {
		initErr = doInit(filePath, name, level)
	})
	return initErr
}

func doInit(filePath, name string, level logrus.Level) error {
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return fmt.Errorf("create log folder %s failed: %v", filePath, err)
	}

	l := logrus.New()
	l.SetLevel(level)
	l.SetReportCaller(false)
	l.SetOutput(os.Stdout) // 控制台仍然输出一份
	l.SetFormatter(&Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		CallerFirst:     true,
		// CustomCallerFormatter: func(f *runtime.Frame) string {
		// 	return fmt.Sprintf(" %s:%d", path.Base(f.File), f.Line)
		// },
	})
	l.AddHook(&FileHook{
		filePath: filePath,
		name:     name,
	})

	std = l
	return nil
}

func SetLevel(level logrus.Level) { ensure().SetLevel(level) }

func Trace(args ...interface{}) { ensure().WithField("caller", caller(2)).Trace(args...) }
func Debug(args ...interface{}) { ensure().WithField("caller", caller(2)).Debug(args...) }
func Info(args ...interface{})  { ensure().WithField("caller", caller(2)).Info(args...) }
func Warn(args ...interface{})  { ensure().WithField("caller", caller(2)).Warn(args...) }
func Error(args ...interface{}) { ensure().WithField("caller", caller(2)).Error(args...) }
func Fatal(args ...interface{}) { ensure().WithField("caller", caller(2)).Fatal(args...) }
func Panic(args ...interface{}) { ensure().WithField("caller", caller(2)).Panic(args...) }

func Traceln(args ...interface{}) { ensure().WithField("caller", caller(2)).Traceln(args...) }
func Debugln(args ...interface{}) { ensure().WithField("caller", caller(2)).Debugln(args...) }
func Infoln(args ...interface{})  { ensure().WithField("caller", caller(2)).Infoln(args...) }
func Warnln(args ...interface{})  { ensure().WithField("caller", caller(2)).Warnln(args...) }
func Errorln(args ...interface{}) { ensure().WithField("caller", caller(2)).Errorln(args...) }
func Fatalln(args ...interface{}) { ensure().WithField("caller", caller(2)).Fatalln(args...) }
func Panicln(args ...interface{}) { ensure().WithField("caller", caller(2)).Panicln(args...) }

func Tracef(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Tracef(format, args...)
}
func Debugf(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	ensure().WithField("caller", caller(2)).Panicf(format, args...)
}

func Close() error {
	if std == nil {
		return nil
	}
	var firstErr error
	for _, hs := range std.Hooks {
		for _, hook := range hs {
			if fh, ok := hook.(*FileHook); ok {
				if err := fh.Close(); err != nil && firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	return firstErr
}

func ensure() *logrus.Logger {
	if std == nil {
		// 未初始化时降级到 stderr，保证不丢日志、不 panic。
		l := logrus.New()
		l.SetOutput(os.Stderr)
		l.SetReportCaller(false)
		l.SetFormatter(&Formatter{
			HideKeys:        true,
			NoColors:        true,
			TimestampFormat: "2006-01-02 15:04:05",
			ShowFullLevel:   true,
			CallerFirst:     true,
		})
		return l
	}
	return std
}

func caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", path.Base(file), line)
}
