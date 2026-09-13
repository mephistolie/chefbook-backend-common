package access

import "github.com/google/uuid"

const (
	ClaimUserId           = "sub"
	ClaimEmail            = "eml"
	ClaimUsername         = "usr"
	ClaimRole             = "rol"
	ClaimSubscriptionPlan = "pln"
	ClaimDeleted          = "del"
	ClaimExpiration       = "exp"
	ClaimNotBefore        = "nbf"
	ClaimIssuedAtTime     = "iat"
)

type Payload struct {
	UserId           uuid.UUID
	Email            string
	Username         *string
	Role             string
	SubscriptionPlan string
	Deleted          bool
}
