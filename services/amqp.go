package services

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPService interface {
	Setup() error
	Send(ctx context.Context, data []byte) error
	GetConsumer() (<-chan amqp.Delivery, error)
}

func NewAMQPService(queueName string) AMQPService {
	return &amqpService{
		queueName: queueName,
	}
}

type amqpService struct {
	queueName string
	channel   *amqp.Channel
}

func (a *amqpService) Setup() error {
	d, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("dial to amqp: %w", err)
	}
	ch, err := d.Channel()
	if err != nil {
		return fmt.Errorf("get channel amqp: %w", err)
	}
	a.channel = ch

	_, err = ch.QueueDeclare(
		a.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declaration: %w", err)
	}

	return nil
}

func (a *amqpService) Send(ctx context.Context, data []byte) error {
	err := a.channel.PublishWithContext(
		ctx,
		"",
		a.queueName,
		false,
		false,
		amqp.Publishing{
			Body: data,
		},
	)

	if err != nil {
		return fmt.Errorf("send to amqp: %w", err)
	}

	return nil
}

func (a *amqpService) GetConsumer() (<-chan amqp.Delivery, error) {
	delivery, err := a.channel.Consume(a.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("send to amqp: %w", err)
	}

	return delivery, nil
}
