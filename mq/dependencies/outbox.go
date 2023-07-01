package dependencies

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
)

type Outbox interface {
	GetPendingMessages() ([]*model.MessageData, error)
	MarkMessageSent(messageId uuid.UUID) error
}
