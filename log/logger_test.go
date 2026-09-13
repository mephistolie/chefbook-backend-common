package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestLogWritesStructuredJSONEvent(t *testing.T) {
	buf := bytes.Buffer{}
	logger = newLogger(&buf, "auth", EnvProd, true)

	Log(context.Background(), Event{
		Event:     "auth.session.created",
		Message:   "session created",
		Component: ComponentGRPC,
		UserID:    "user-id",
		Duration:  24 * time.Millisecond,
		Payload: map[string]any{
			"client": "mobile",
		},
	})

	entry := decodeLogEntry(t, buf.Bytes())
	assertLogField(t, entry, FieldService, "auth")
	assertLogField(t, entry, FieldEnvironment, string(EnvProd))
	assertLogField(t, entry, FieldEvent, "auth.session.created")
	assertLogField(t, entry, "message", "session created")
	assertLogField(t, entry, FieldComponent, ComponentGRPC)
	assertLogField(t, entry, FieldUserID, "user-id")
	assertLogField(t, entry, FieldDurationMs, float64(24))

	payload, ok := entry[FieldPayload].(map[string]any)
	if !ok {
		t.Fatalf("expected payload object, got %#v", entry[FieldPayload])
	}
	assertLogField(t, payload, "client", "mobile")
}

func TestLogErrorWritesErrorField(t *testing.T) {
	buf := bytes.Buffer{}
	logger = newLogger(&buf, "recipe", EnvDev, true)

	LogError(context.Background(), Event{
		Event:     "postgres.query.failed",
		Message:   "unable to get recipe",
		Component: ComponentPostgres,
		ErrorType: "not_found",
	}, errors.New("query failed"))

	entry := decodeLogEntry(t, buf.Bytes())
	assertLogField(t, entry, FieldService, "recipe")
	assertLogField(t, entry, FieldEnvironment, string(EnvDev))
	assertLogField(t, entry, FieldEvent, "postgres.query.failed")
	assertLogField(t, entry, FieldErrorType, "not_found")
	assertLogField(t, entry, FieldError, "query failed")
}

func decodeLogEntry(t *testing.T, data []byte) map[string]any {
	t.Helper()

	var entry map[string]any
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("unable to decode log entry %s: %v", string(data), err)
	}
	return entry
}

func assertLogField(t *testing.T, entry map[string]any, key string, want any) {
	t.Helper()

	if got := entry[key]; got != want {
		t.Fatalf("expected %s=%#v, got %#v in %#v", key, want, got, entry)
	}
}
