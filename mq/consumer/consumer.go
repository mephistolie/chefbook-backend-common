package consumer

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/config"
	"github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	amqp "github.com/wagslane/go-rabbitmq"
	"k8s.io/utils/strings/slices"
)

type Consumer struct {
	conn              *amqp.Conn
	consumers         []*amqp.Consumer
	inbox             dependencies.Inbox
	supportedMsgTypes []string
}

type Params struct {
	QueueName string
	Options   []func(*amqp.ConsumerOptions)
}

func New(cfg config.Amqp, inbox dependencies.Inbox, supportedMsgTypes []string) (*Consumer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", *cfg.User, *cfg.Password, *cfg.Host, *cfg.Port, *cfg.VHost)
	conn, err := amqp.NewConn(url)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:              conn,
		inbox:             inbox,
		supportedMsgTypes: supportedMsgTypes,
	}, nil
}

func (s *Consumer) Start(consumers ...Params) error {
	for _, consumerParams := range consumers {
		consumer, err := amqp.NewConsumer(
			s.conn,
			s.handleDelivery,
			consumerParams.QueueName,
			consumerParams.Options...,
		)
		if err != nil {
			return err
		}
		s.consumers = append(s.consumers, consumer)
	}

	return nil
}

func (s *Consumer) handleDelivery(delivery amqp.Delivery) amqp.Action {
	messageId, err := uuid.Parse(delivery.MessageId)
	if err != nil {
		log.Warn("invalid message id: ", delivery.MessageId)
		return amqp.NackDiscard
	}

	if !slices.Contains(s.supportedMsgTypes, delivery.Type) {
		log.Infof("unsupported message type: ", delivery.Type)
		return amqp.NackDiscard
	}

	msg := model.MessageData{
		Id:       messageId,
		Exchange: delivery.Exchange,
		Type:     delivery.Type,
		Body:     delivery.Body,
	}
	if err = s.inbox.HandleMessage(msg); err != nil {
		log.Warn("requeue message ", msg.Id)
		return amqp.NackRequeue
	}

	return amqp.Ack
}

func (s *Consumer) Stop() error {
	for _, consumer := range s.consumers {
		consumer.Close()
	}
	return s.conn.Close()
}
