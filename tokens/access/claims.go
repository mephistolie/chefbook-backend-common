package access

const (
	ClaimUserId       = "sub"
	ClaimEmail        = "eml"
	ClaimNickname     = "nik"
	ClaimRole         = "rol"
	ClaimPremium      = "prm"
	ClaimExpiration   = "exp"
	ClaimNotBefore    = "nbf"
	ClaimIssuedAtTime = "iat"
)

type Payload struct {
	UserId   string
	Email    string
	Nickname *string
	Role     string
	Premium  bool
}
