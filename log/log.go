package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const StackTraceSkip = 2

var Logger = logrus.New()

type Fields logrus.Fields

func SetLogLevel(level logrus.Level) {
	Logger.Level = level
}

func SetLogFormatter(formatter logrus.Formatter) {
	Logger.Formatter = formatter
}

func Info(args ...interface{}) {
	if Logger.Level >= logrus.InfoLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Info(args...)
	}
}

func InfoWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.InfoLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Info(l)
	}
}

func Debug(args ...interface{}) {
	if Logger.Level >= logrus.DebugLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Debug(args...)
	}
}

func DebugWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.DebugLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Debug(l)
	}
}

func Warn(args ...interface{}) {
	if Logger.Level >= logrus.WarnLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Warn(args...)
	}
}

func WarnWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.WarnLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Warn(l)
	}
}

func Error(args ...interface{}) {
	if Logger.Level >= logrus.ErrorLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Error(args...)
	}
}

func ErrorWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.ErrorLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Error(l)
	}
}

func Fatal(args ...interface{}) {
	if Logger.Level >= logrus.FatalLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Fatal(args...)
	}
}

func FatalWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.FatalLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Fatal(l)
	}
}

func Panic(args ...interface{}) {
	if Logger.Level >= logrus.PanicLevel {
		logger := Logger.WithFields(logrus.Fields{})
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Panic(args...)
	}
}

func PanicWithFields(l interface{}, f Fields) {
	if Logger.Level >= logrus.PanicLevel {
		logger := Logger.WithFields(logrus.Fields(f))
		logger.Data["file"] = getFile(StackTraceSkip)
		logger.Panic(l)
	}
}

func WithFields(fields logrus.Fields) {
	Logger.WithFields(fields)
}

func getFile(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	slash := strings.LastIndex(file, "/")
	if slash >= 0 {
		preSlash := strings.LastIndex(file[0:slash], "/")
		file = file[preSlash+1:slash] + file[slash:]
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
}
