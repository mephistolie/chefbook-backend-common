package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/config"
	outbox "github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	amqp "github.com/wagslane/go-rabbitmq"
)

type Publisher struct {
	appId             string
	conn              *amqp.Conn
	publisherProfiles *amqp.Publisher
	outbox            outbox.Outbox
}

func New(appId string, cfg config.Amqp, outbox outbox.Outbox) (*Publisher, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", *cfg.User, *cfg.Password, *cfg.Host, *cfg.Port, *cfg.VHost)
	conn, err := amqp.NewConn(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		appId:  appId,
		conn:   conn,
		outbox: outbox,
	}, nil
}

func (p *Publisher) Start(publisherOptions ...func(*amqp.PublisherOptions)) error {
	var err error = nil
	p.publisherProfiles, err = amqp.NewPublisher(
		p.conn,
		publisherOptions...,
	)
	if err != nil {
		return err
	}

	go p.observeOutbox()

	return nil
}

func (p *Publisher) observeOutbox() {
	for {
		fails := 0
		if msgs, err := p.outbox.GetPendingMessages(); err == nil {
			for _, msg := range msgs {
				if err = p.PublishMessage(msg); err != nil {
					fails += 1
					if fails >= 5 {
						break
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (p *Publisher) PublishMessage(msg *model.MessageData) error {
	log.Log(context.Background(), log.Event{
		Event:     "mq.message.publishing",
		Message:   "publishing message",
		Component: log.ComponentAMQP,
		MessageID: msg.Id.String(),
		Payload: map[string]any{
			"message_type": msg.Type,
			"exchange":     msg.Exchange,
		},
	})
	err := p.publisherProfiles.Publish(
		msg.Body,
		[]string{""},
		amqp.WithPublishOptionsExchange(msg.Exchange),
		amqp.WithPublishOptionsMessageID(msg.Id.String()),
		amqp.WithPublishOptionsPersistentDelivery,
		amqp.WithPublishOptionsContentType("application/json"),
		amqp.WithPublishOptionsType(msg.Type),
		amqp.WithPublishOptionsAppID(p.appId),
	)
	if err == nil {
		log.Log(context.Background(), log.Event{
			Event:     "mq.message.published",
			Message:   "message published",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
				"exchange":     msg.Exchange,
			},
		})
	} else {
		log.LogWarnError(context.Background(), log.Event{
			Event:     "mq.message.publish_failed",
			Message:   "unable to publish message",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
				"exchange":     msg.Exchange,
			},
		}, err)
	}

	if err == nil {
		_ = p.outbox.MarkMessageSent(msg.Id)
	}

	return err
}

func (p *Publisher) Stop() error {
	p.publisherProfiles.Close()
	return p.conn.Close()
}
