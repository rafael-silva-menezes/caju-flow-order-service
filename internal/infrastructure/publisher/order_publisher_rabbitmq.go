package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"order-service/internal/domain/entity"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
	routingKey string
}

func NewRabbitMQPublisher(conn *amqp091.Connection, exchange, routingKey, queueName string) (*RabbitMQPublisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = channel.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return &RabbitMQPublisher{
		connection: conn,
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

func (p *RabbitMQPublisher) PublishCreatedOrder(ctx context.Context, order *entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Order %s published to RabbitMQ", order.ID)
	return nil
}

func (p *RabbitMQPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := p.connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
