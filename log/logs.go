package log

import (
	"context"

	"github.com/rs/zerolog"
)

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
