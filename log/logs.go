package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
)

func Trace(msg ...interface{}) {
	callerLogger().Trace(msg...)
}

func Tracef(format string, args ...interface{}) {
	callerLogger().Tracef(format, args...)
}

func Debug(msg ...interface{}) {
	callerLogger().Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	callerLogger().Debugf(format, args...)
}

func Info(msg ...interface{}) {
	callerLogger().Info(msg...)
}

func Infof(format string, args ...interface{}) {
	callerLogger().Infof(format, args...)
}

func Warn(msg ...interface{}) {
	callerLogger().Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	callerLogger().Warnf(format, args...)
}

func Error(msg ...interface{}) {
	callerLogger().Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	callerLogger().Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	callerLogger().Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	callerLogger().Fatalf(format, args...)
}

func callerLogger() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return e
	}
	fun := runtime.FuncForPC(pc)

	fields := logrus.Fields{}
	fields["file"] = fmt.Sprintf("%s:%d", path.Base(file), line)
	fields["func"] = fmt.Sprintf("%s()", fun.Name())

	return e.WithFields(fields)
}
