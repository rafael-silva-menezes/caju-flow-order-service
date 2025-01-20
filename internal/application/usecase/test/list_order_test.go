package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/application/usecase"
	usecasemock "order-service/internal/application/usecase/mock"
	"order-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListOrderUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		useCase := usecase.NewListOrderUseCase(mockRepo)

		items := []entity.Item{
			{
				ID:       "1",
				Name:     "Item 1",
				Quantity: 2,
				Price:    10.0,
			},
		}

		order1, _ := entity.NewOrder("order123", "John Doe", items)
		order2, _ := entity.NewOrder("order456", "Jane Doe", items)
		orders := []entity.Order{*order1, *order2}

		mockRepo.On("List", mock.Anything).Return(orders, nil)

		output, err := useCase.Execute(context.Background())

		assert.NoError(t, err)
		assert.Len(t, output.Orders, 2)
		assert.Equal(t, order1.OrderID, output.Orders[0].OrderID)
		assert.Equal(t, order2.OrderID, output.Orders[1].OrderID)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		useCase := usecase.NewListOrderUseCase(mockRepo)

		mockRepo.On("List", mock.Anything).Return([]entity.Order{}, nil)

		output, err := useCase.Execute(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, output.Orders)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		useCase := usecase.NewListOrderUseCase(mockRepo)

		mockRepo.On("List", mock.Anything).Return([]entity.Order{}, errors.New("database error"))

		output, err := useCase.Execute(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Empty(t, output.Orders)
	})
}
