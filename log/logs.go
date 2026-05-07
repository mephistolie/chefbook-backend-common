package log

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/rs/zerolog"
)

func Trace(msg ...interface{}) {
	legacyEvent(logger.Trace(), msg...)
}

func Tracef(format string, args ...interface{}) {
	legacyEventf(logger.Trace(), format, args...)
}

func Debug(msg ...interface{}) {
	legacyEvent(logger.Debug(), msg...)
}

func Debugf(format string, args ...interface{}) {
	legacyEventf(logger.Debug(), format, args...)
}

func Info(msg ...interface{}) {
	legacyEvent(logger.Info(), msg...)
}

func Infof(format string, args ...interface{}) {
	legacyEventf(logger.Info(), format, args...)
}

func Warn(msg ...interface{}) {
	legacyEvent(logger.Warn(), msg...)
}

func Warnf(format string, args ...interface{}) {
	legacyEventf(logger.Warn(), format, args...)
}

func Error(msg ...interface{}) {
	legacyEvent(logger.Error(), msg...)
}

func Errorf(format string, args ...interface{}) {
	legacyEventf(logger.Error(), format, args...)
}

func Fatal(msg ...interface{}) {
	legacyEvent(logger.Fatal(), msg...)
}

func Fatalf(format string, args ...interface{}) {
	legacyEventf(logger.Fatal(), format, args...)
}

func Log(ctx context.Context, event Event) {
	writeEvent(logger.Info().Ctx(ctx), event, nil)
}

func LogDebug(ctx context.Context, event Event) {
	writeEvent(logger.Debug().Ctx(ctx), event, nil)
}

func LogWarn(ctx context.Context, event Event) {
	writeEvent(logger.Warn().Ctx(ctx), event, nil)
}

func LogWarnError(ctx context.Context, event Event, err error) {
	writeEvent(logger.Warn().Ctx(ctx), event, err)
}

func LogError(ctx context.Context, event Event, err error) {
	writeEvent(logger.Error().Ctx(ctx), event, err)
}

func LogFatal(ctx context.Context, event Event, err error) {
	writeEvent(logger.Fatal().Ctx(ctx), event, err)
}

func legacyEvent(event *zerolog.Event, msg ...interface{}) {
	event.
		Str(FieldEvent, callerEventName()).
		Msg(fmt.Sprint(msg...))
}

func legacyEventf(event *zerolog.Event, format string, args ...interface{}) {
	event.
		Str(FieldEvent, callerEventName()).
		Msgf(format, args...)
}

func callerEventName() string {
	pc, file, _, ok := runtime.Caller(3)
	if !ok {
		return "log.callsite.unknown"
	}

	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = fn.Name()
		if index := strings.LastIndex(funcName, "/"); index >= 0 {
			funcName = funcName[index+1:]
		}
		funcName = strings.TrimPrefix(funcName, serviceName+".")
	}

	component := componentFromPath(file)
	if component == "" {
		component = "app"
	}

	return strings.Join([]string{
		normalizeEventPart(serviceName),
		normalizeEventPart(component),
		normalizeEventPart(funcName),
	}, ".")
}

func componentFromPath(file string) string {
	path := filepath.ToSlash(file)
	switch {
	case strings.Contains(path, "/repository/postgres/"):
		return ComponentPostgres
	case strings.Contains(path, "/repository/s3/"):
		return ComponentS3
	case strings.Contains(path, "/amqp/"), strings.Contains(path, "/mq/"):
		return ComponentAMQP
	case strings.Contains(path, "/grpc/"):
		return ComponentGRPC
	case strings.Contains(path, "/http/"):
		return ComponentHTTP
	case strings.Contains(path, "/config/"):
		return "config"
	case strings.Contains(path, "/app/"):
		return "app"
	default:
		return filepath.Base(filepath.Dir(path))
	}
}

func normalizeEventPart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}

	var builder strings.Builder
	var previousUnderscore bool
	for _, r := range value {
		switch {
		case r == '.' || r == '-' || r == '/' || r == '*' || r == '(' || r == ')':
			if !previousUnderscore {
				builder.WriteByte('_')
				previousUnderscore = true
			}
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if unicode.IsUpper(r) && builder.Len() > 0 && !previousUnderscore {
				builder.WriteByte('_')
			}
			builder.WriteRune(unicode.ToLower(r))
			previousUnderscore = false
		default:
			if !previousUnderscore {
				builder.WriteByte('_')
				previousUnderscore = true
			}
		}
	}

	result := strings.Trim(builder.String(), "_")
	if result == "" {
		return "unknown"
	}
	return result
}

func writeEvent(logEvent *zerolog.Event, event Event, err error) {
	if event.Event == "" {
		event.Event = "unknown"
	}
	if event.Message == "" {
		event.Message = event.Event
	}

	logEvent = logEvent.Str(FieldEvent, event.Event)
	if event.Component != "" {
		logEvent = logEvent.Str(FieldComponent, event.Component)
	}
	if event.RequestID != "" {
		logEvent = logEvent.Str(FieldRequestID, event.RequestID)
	}
	if event.TraceID != "" {
		logEvent = logEvent.Str(FieldTraceID, event.TraceID)
	}
	if event.SpanID != "" {
		logEvent = logEvent.Str(FieldSpanID, event.SpanID)
	}
	if event.UserID != "" {
		logEvent = logEvent.Str(FieldUserID, event.UserID)
	}
	if event.MessageID != "" {
		logEvent = logEvent.Str(FieldMessageID, event.MessageID)
	}
	if event.Operation != "" {
		logEvent = logEvent.Str(FieldOperation, event.Operation)
	}
	if event.Duration > 0 {
		logEvent = logEvent.Int64(FieldDurationMs, event.Duration.Milliseconds())
	}
	if event.ErrorType != "" {
		logEvent = logEvent.Str(FieldErrorType, event.ErrorType)
	}
	if event.GRPCMethod != "" {
		logEvent = logEvent.Str(FieldGRPCMethod, event.GRPCMethod)
	}
	if event.GRPCCode != "" {
		logEvent = logEvent.Str(FieldGRPCCode, event.GRPCCode)
	}
	if event.HTTPMethod != "" {
		logEvent = logEvent.Str(FieldHTTPMethod, event.HTTPMethod)
	}
	if event.HTTPPath != "" {
		logEvent = logEvent.Str(FieldHTTPPath, event.HTTPPath)
	}
	if event.HTTPStatus > 0 {
		logEvent = logEvent.Int(FieldHTTPStatus, event.HTTPStatus)
	}
	if event.Payload != nil {
		logEvent = logEvent.Interface(FieldPayload, event.Payload)
	}
	if err != nil {
		logEvent = logEvent.Err(err)
	}
	logEvent.Msg(event.Message)
}
