package dependencies

import (
	"github.com/mephistolie/chefbook-backend-common/mq/model"
)

type Inbox interface {
	HandleMessage(msg model.MessageData) error
}
