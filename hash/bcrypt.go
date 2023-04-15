package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptManager struct {
	saltCost int
}

func NewBcryptManager(saltCost int) *BcryptManager {
	return &BcryptManager{saltCost: saltCost}
}

func (m *BcryptManager) Hash(data string) (string, error) {
	hashedData, err := bcrypt.GenerateFromPassword([]byte(data), m.saltCost)
	return string(hashedData), err
}

func (m *BcryptManager) Validate(data string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
}
