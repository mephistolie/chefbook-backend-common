package fail

const (
	TypeUnknown         = "unknown"
	TypeInvalidBody     = "invalid_body"
	TypeUnauthorized    = "unauthorized"
	TypeAccessDenied    = "access_denied"
	TypeNotFound        = "not_found"
	TypeConflict        = "conflict"
	TypePremiumRequired = "premium_required"
	TypeUnavailable     = "unavailable"
)

type Response struct {
	Code      int    `json:"-"`
	ErrorType string `json:"error"`
	Message   string `json:"message,omitempty"`
}
