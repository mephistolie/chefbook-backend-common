package log

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

func Trace(msg ...interface{}) {
	legacyEvent(logger.Trace(), msg...)
}

func Tracef(format string, args ...interface{}) {
	logger.Trace().Msgf(format, args...)
}

func Debug(msg ...interface{}) {
	legacyEvent(logger.Debug(), msg...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debug().Msgf(format, args...)
}

func Info(msg ...interface{}) {
	legacyEvent(logger.Info(), msg...)
}

func Infof(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

func Warn(msg ...interface{}) {
	legacyEvent(logger.Warn(), msg...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warn().Msgf(format, args...)
}

func Error(msg ...interface{}) {
	legacyEvent(logger.Error(), msg...)
}

func Errorf(format string, args ...interface{}) {
	logger.Error().Msgf(format, args...)
}

func Fatal(msg ...interface{}) {
	legacyEvent(logger.Fatal(), msg...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal().Msgf(format, args...)
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

func LogError(ctx context.Context, event Event, err error) {
	writeEvent(logger.Error().Ctx(ctx), event, err)
}

func legacyEvent(event *zerolog.Event, msg ...interface{}) {
	event.
		Str(FieldEvent, "legacy.log").
		Msg(fmt.Sprint(msg...))
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
