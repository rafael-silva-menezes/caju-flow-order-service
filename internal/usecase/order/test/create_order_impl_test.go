package order_test

import (
	"context"
	"testing"

	"order-service/internal/usecase/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully creates an order", func(t *testing.T) {
		mockRepo := new(mockOrderRepository)
		mockPub := new(mockOrderPublisher)
		useCase := order.NewCreateOrderUseCase(mockRepo, mockPub)

		input := order.CreateOrderInput{
			CustomerName: "John Doe",
			Items: []order.ItemInput{
				{ID: "1", Name: "Item1", Quantity: 2, Price: 10.0},
			},
		}

		mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
		mockPub.On("PublishCreatedOrder", mock.Anything, mock.Anything).Return(nil)

		output, err := useCase.Execute(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, "John Doe", output.CustomerName)
		assert.Equal(t, 20.0, output.Total)
		mockRepo.AssertExpectations(t)
		mockPub.AssertExpectations(t)
	})

	t.Run("fails when order has no items", func(t *testing.T) {
		mockRepo := new(mockOrderRepository)
		mockPub := new(mockOrderPublisher)
		useCase := order.NewCreateOrderUseCase(mockRepo, mockPub)

		input := order.CreateOrderInput{
			CustomerName: "John Doe",
			Items:        []order.ItemInput{},
		}

		_, err := useCase.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, "order must contain at least one item", err.Error())
	})
}
