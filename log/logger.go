package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)

var (
	e *logrus.Entry
)

func Init(
	logsPath string,
	debug bool,
) {
	l := newLogger()

	writer := []io.Writer{os.Stdout}
	if len(logsPath) > 0 {
		err := os.MkdirAll(path.Dir(logsPath), 0755)
		if err != nil || os.IsExist(err) {
			panic("can't create log dir")
		}

		logFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("can't open logs file: %s", err))
		}
		writer = []io.Writer{logFile, os.Stdout}
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    writer,
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(getLogLevel(debug))

	e = logrus.NewEntry(l)
}

func newLogger() *logrus.Logger {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	return l
}

func getLogLevel(debug bool) logrus.Level {
	logLevel := logrus.InfoLevel
	if debug {
		logLevel = logrus.TraceLevel
	}
	return logLevel
}
