package access

import "github.com/google/uuid"

const (
	ClaimUserId           = "sub"
	ClaimEmail            = "eml"
	ClaimNickname         = "nik"
	ClaimRole             = "rol"
	ClaimSubscriptionPlan = "pln"
	ClaimExpiration       = "exp"
	ClaimNotBefore        = "nbf"
	ClaimIssuedAtTime     = "iat"
)

type Payload struct {
	UserId           uuid.UUID
	Email            string
	Nickname         *string
	Role             string
	SubscriptionPlan string
}
