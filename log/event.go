package log

import "time"

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvProd Environment = "prod"
)

const (
	ComponentAMQP      = "amqp"
	ComponentFirebase  = "firebase"
	ComponentGRPC      = "grpc"
	ComponentHTTP      = "http"
	ComponentPostgres  = "postgres"
	ComponentS3        = "s3"
	ComponentScheduler = "scheduler"
)

const (
	FieldComponent   = "component"
	FieldDurationMs  = "duration_ms"
	FieldEnvironment = "environment"
	FieldError       = "error"
	FieldErrorType   = "error_type"
	FieldEvent       = "event"
	FieldGRPCCode    = "grpc_code"
	FieldGRPCMethod  = "grpc_method"
	FieldHTTPMethod  = "http_method"
	FieldHTTPPath    = "http_path"
	FieldHTTPStatus  = "http_status"
	FieldMessageID   = "message_id"
	FieldOperation   = "operation"
	FieldPayload     = "payload"
	FieldRequestID   = "request_id"
	FieldService     = "service"
	FieldSpanID      = "span_id"
	FieldTraceID     = "trace_id"
	FieldUserID      = "user_id"
)

type Event struct {
	Event     string
	Message   string
	Component string

	RequestID string
	TraceID   string
	SpanID    string
	UserID    string
	MessageID string
	Operation string

	Duration  time.Duration
	ErrorType string

	GRPCMethod string
	GRPCCode   string

	HTTPMethod string
	HTTPPath   string
	HTTPStatus int

	Payload map[string]any
}
