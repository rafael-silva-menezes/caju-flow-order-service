package publisher_test

import (
	"context"
	"testing"
	"time"

	"order-service/internal/domain/entity"
	"order-service/internal/infrastructure/publisher"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestPublishCreatedOrder_Integration(t *testing.T) {
	conn, err := amqp091.Dial("amqp://order_test_user:order_test_pass@localhost:5675/")
	assert.NoError(t, err)
	defer conn.Close()

	orderPublisher, err := publisher.NewRabbitMQPublisher(conn, "order-exchange", "order-routing-key")
	assert.NoError(t, err)

	items := []entity.Item{
		{ID: "item1", Quantity: 1, Price: 100.0},
	}

	order, err := entity.NewOrder("123", "Customer A", items)
	assert.NoError(t, err)

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"order-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)

	err = ch.ExchangeDeclare(
		"order-exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)

	err = ch.QueueBind(
		"order-queue",
		"order-routing-key",
		"order-exchange",
		false,
		nil,
	)
	assert.NoError(t, err)

	msgs, err := ch.Consume(
		"order-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)

	err = orderPublisher.PublishCreatedOrder(context.Background(), order)
	assert.NoError(t, err)

	select {
	case msg := <-msgs:
		assert.Equal(t, "application/json", msg.ContentType)
		assert.Contains(t, string(msg.Body), `"id":"123"`)
	case <-time.After(1 * time.Second):
		t.Fatal("Mensagem não recebida no tempo esperado")
	}
}
