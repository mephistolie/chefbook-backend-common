package consumer

import (
	"context"
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
			consumerParams.QueueName,
			consumerParams.Options...,
		)
		if err != nil {
			return err
		}
		s.consumers = append(s.consumers, consumer)
		go func(consumer *amqp.Consumer) {
			if err := consumer.Run(s.handleDelivery); err != nil {
				log.LogWarnError(context.Background(), log.Event{
					Event:     "mq.consumer.stopped",
					Message:   "rabbitmq consumer stopped with error",
					Component: log.ComponentAMQP,
				}, err)
			}
		}(consumer)
	}

	return nil
}

func (s *Consumer) handleDelivery(delivery amqp.Delivery) amqp.Action {
	messageId, err := uuid.Parse(delivery.MessageId)
	if err != nil {
		log.LogWarn(context.Background(), log.Event{
			Event:     "mq.message.invalid_id",
			Message:   "invalid message id",
			Component: log.ComponentAMQP,
			Payload: map[string]any{
				"raw_message_id": delivery.MessageId,
				"message_type":   delivery.Type,
			},
		})
		return amqp.NackDiscard
	}

	if !slices.Contains(s.supportedMsgTypes, delivery.Type) {
		log.LogWarn(context.Background(), log.Event{
			Event:     "mq.message.unsupported_type",
			Message:   "unsupported message type",
			Component: log.ComponentAMQP,
			MessageID: messageId.String(),
			Payload: map[string]any{
				"message_type": delivery.Type,
				"exchange":     delivery.Exchange,
			},
		})
		return amqp.NackDiscard
	}

	msg := model.MessageData{
		Id:       messageId,
		Exchange: delivery.Exchange,
		Type:     delivery.Type,
		Body:     delivery.Body,
	}
	if err = s.inbox.HandleMessage(msg); err != nil {
		log.LogWarnError(context.Background(), log.Event{
			Event:     "mq.message.requeued",
			Message:   "message requeued",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
				"exchange":     msg.Exchange,
			},
		}, err)
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
