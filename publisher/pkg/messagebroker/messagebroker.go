package messagebroker

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(theme string, data any) error
}

type EventPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewEventPublisher(url string) (*EventPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connecting to messagebroker: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("opening channel: %w", err)
	}

	return &EventPublisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p EventPublisher) Publish(topic string, data any) error {
	err := p.ch.ExchangeDeclare(
		topic,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declaring exchange: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encoding data: %w", err)
	}

	err = p.ch.Publish(
		topic,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		})
	if err != nil {
		return fmt.Errorf("declaring exchange: %w", err)
	}

	return nil
}
