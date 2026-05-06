package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger      zerolog.Logger
	serviceName = "unknown"
	environment = EnvProd
)

func init() {
	logger = newLogger(os.Stdout, serviceName, environment, false)
}

func Init(
	logsPath string,
	debug bool,
) {
	InitWithService("unknown", logsPath, debug)
}

func InitWithService(
	service string,
	logsPath string,
	debug bool,
) {
	serviceName = service
	environment = getEnvironment(debug)

	writer := io.Writer(os.Stdout)
	if len(logsPath) > 0 {
		err := os.MkdirAll(path.Dir(logsPath), 0755)
		if err != nil || os.IsExist(err) {
			panic("can't create log dir")
		}

		logFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("can't open logs file: %s", err))
		}
		writer = zerolog.MultiLevelWriter(logFile, os.Stdout)
	}

	logger = newLogger(writer, serviceName, environment, debug)
}

func newLogger(writer io.Writer, service string, env Environment, debug bool) zerolog.Logger {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = FieldError
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(getLogLevel(debug))

	return zerolog.New(writer).
		With().
		Timestamp().
		Str(FieldService, service).
		Str(FieldEnvironment, string(env)).
		Logger()
}

func getLogLevel(debug bool) zerolog.Level {
	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.TraceLevel
	}
	return logLevel
}

func getEnvironment(debug bool) Environment {
	if debug {
		return EnvDev
	}
	return EnvProd
}
