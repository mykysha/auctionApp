package messagebroker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber interface {
	SubscribeToTopic(topic string) error
}

type EventSubscriber struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	msgChan chan []byte
}

func NewEventSubscriber(url string, msgChan chan []byte) (*EventSubscriber, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connecting to messagebroker: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("opening channel: %w", err)
	}

	return &EventSubscriber{
		conn:    conn,
		ch:      ch,
		msgChan: msgChan,
	}, nil
}

func (s EventSubscriber) SubscribeToTopic(topic string) error {
	err := s.ch.ExchangeDeclare(
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

	queue, err := s.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declaring queue: %w", err)
	}

	err = s.ch.QueueBind(
		queue.Name,
		"",
		topic,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("binding queue: %w", err)
	}

	messages, err := s.ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("registering consumer: %w", err)
	}

	go func() {
		for delivery := range messages {
			s.msgChan <- delivery.Body
		}
	}()

	return nil
}
