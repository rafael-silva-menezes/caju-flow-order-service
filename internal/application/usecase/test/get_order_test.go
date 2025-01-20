package usecase_test

import (
	"context"
	"testing"

	"order-service/internal/application/dtos"
	"order-service/internal/application/usecase"
	usecasemock "order-service/internal/application/usecase/mock"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrderUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		useCase := usecase.NewGetOrderUseCase(mockRepo)

		items := []entity.Item{
			{
				ID:       "1",
				Name:     "Item 1",
				Quantity: 2,
				Price:    10.0,
			},
		}

		order, _ := entity.NewOrder("order123", "John Doe", items)
		mockRepo.On("FindByID", mock.Anything, "order123").Return(order, nil)

		expectedOutput := dtos.FromEntityToGetOrderOutput(order)

		output, err := useCase.Execute(context.Background(), "order123")

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput.OrderID, output.OrderID)
		assert.Equal(t, expectedOutput.CustomerName, output.CustomerName)
		assert.Equal(t, expectedOutput.Total, output.Total)
		assert.Equal(t, expectedOutput.Status, output.Status)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		useCase := usecase.NewGetOrderUseCase(mockRepo)

		mockRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, repository.ErrNotFound)

		output, err := useCase.Execute(context.Background(), "non-existent")

		assert.Error(t, err)
		assert.Equal(t, "order not found", err.Error())
		assert.Empty(t, output)
	})
}
