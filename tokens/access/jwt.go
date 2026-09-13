package access

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

type Producer struct {
	key *rsa.PrivateKey
}

func NewProducer(privateKeyPath string) (*Producer, *rsa.PublicKey, error) {
	data, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, nil, err
	}

	return &Producer{key: key}, &key.PublicKey, nil
}

func NewProducerByRawKey(data []byte) (*Producer, *rsa.PublicKey, error) {
	key, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		key, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	}
	if err != nil {
		return nil, nil, err
	}

	return &Producer{key: key}, &key.PublicKey, nil
}

func NewProducerByKey(key *rsa.PrivateKey) *Producer {
	return &Producer{key: key}
}

func (p *Producer) Produce(payload Payload, ttl time.Duration) (string, error) {

	currentTime := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims[ClaimUserId] = payload.UserId.String()
	claims[ClaimEmail] = payload.Email
	if payload.Username != nil {
		claims[ClaimUsername] = *payload.Username
	}
	claims[ClaimRole] = payload.Role
	claims[ClaimSubscriptionPlan] = payload.SubscriptionPlan
	claims[ClaimDeleted] = payload.Deleted
	claims[ClaimExpiration] = currentTime.Add(ttl).Unix()
	claims[ClaimIssuedAtTime] = currentTime.Unix()
	claims[ClaimNotBefore] = currentTime.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(p.key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

type Parser struct {
	Key *rsa.PublicKey
}

func NewParser(publicKeyPath string) (*Parser, error) {
	data, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, err
	}

	return &Parser{Key: key}, nil
}

func NewParserByRawKey(data []byte) (*Parser, error) {
	key, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		key, err = jwt.ParseRSAPublicKeyFromPEM(data)
	}
	if err != nil {
		return nil, err
	}

	return &Parser{Key: key}, nil
}

func NewParserByKey(key *rsa.PublicKey) *Parser {
	return &Parser{Key: key}
}

func (p *Parser) Parse(token string) (Payload, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		return p.Key, nil
	})
	if err != nil {
		return Payload{}, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return Payload{}, fmt.Errorf("invalid token")
	}

	userIdStr, _ := claims[ClaimUserId].(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return Payload{}, err
	}
	email, _ := claims[ClaimEmail].(string)
	var usernamePtr *string = nil
	if username, ok := claims[ClaimUsername].(string); ok && len(username) > 0 {
		usernamePtr = &username
	}
	role, _ := claims[ClaimRole].(string)
	plan, _ := claims[ClaimSubscriptionPlan].(string)
	deleted, _ := claims[ClaimDeleted].(bool)

	return Payload{
		UserId:           userId,
		Email:            email,
		Username:         usernamePtr,
		Role:             role,
		SubscriptionPlan: plan,
		Deleted:          deleted,
	}, nil
}
