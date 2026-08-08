package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "jobs"

type Publisher struct {
	conn *amqp.Connection
}

func NewPublisher() (*Publisher, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	return &Publisher{conn: conn}, nil
}

func (p *Publisher) Close() error {
	return p.conn.Close()
}

func (p *Publisher) PublishPostCreated(
	ctx context.Context,
	message PostCreatedMessage,
) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal post-created message: %w", err)
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("open RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	err = ch.PublishWithContext(
		ctx,
		"",        // default exchange
		queueName, // routes directly to the "jobs" queue
		false,     // fail if no matching queue instead of returning it
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Type:         PostCreated,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish post-created message: %w", err)
	}

	return nil
}
