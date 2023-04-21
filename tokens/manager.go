package tokens

import (
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/tokens/access"
	"time"
)

type Manager struct {
	jwtProducer access.Producer
	jwtParser   access.Parser
}

func NewManager(privateKeyPath, publicKeyPath string) (*Manager, error) {
	jwtProducer, err := access.NewProducer(privateKeyPath)
	if err != nil {
		return nil, err
	}
	jwtParser, err := access.NewParser(publicKeyPath)
	if err != nil {
		return nil, err
	}

	return &Manager{jwtProducer: *jwtProducer, jwtParser: *jwtParser}, nil
}

func NewManagerByKeys(privateKey, publicKey []byte) (*Manager, error) {
	jwtProducer, err := access.NewProducerByKey(privateKey)
	if err != nil {
		return nil, err
	}
	jwtParser, err := access.NewParserByKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &Manager{jwtProducer: *jwtProducer, jwtParser: *jwtParser}, nil
}

func (m *Manager) CreateAccess(payload access.Payload, ttl time.Duration) (string, error) {
	return m.jwtProducer.Produce(payload, ttl)
}

func (m *Manager) ParseAccess(token string) (access.Payload, error) {
	return m.jwtParser.Parse(token)
}

func (m *Manager) GetAccessPublicKey() *rsa.PublicKey {
	return m.jwtParser.Key
}

func (m *Manager) CreateRefresh() string {
	return uuid.New().String()
}
