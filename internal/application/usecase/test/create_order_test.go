package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/application/dtos"
	"order-service/internal/application/usecase"
	usecasemock "order-service/internal/application/usecase/mock"
	"order-service/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderUseCase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		mockPub := new(usecasemock.MockOrderPublisher)
		useCase := usecase.NewCreateOrderUseCase(mockRepo, mockPub)

		input := dtos.OrderInput{
			CustomerName: "John Doe",
			Items: []dtos.ItemInput{
				{
					ID:       "1",
					Name:     "Item 1",
					Quantity: 2,
					Price:    10.0,
				},
			},
		}

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)
		mockPub.On("PublishCreatedOrder", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		mockOrder := &entity.Order{
			ID:           uuid.New().String(),
			CustomerName: "John Doe",
			Status:       entity.Pending,
			Items:        []entity.Item{{ID: "1", Name: "Item 1", Quantity: 2, Price: 10.0}},
		}

		expectedOutput := dtos.FromEntityToOrderOutput(mockOrder)

		output, err := useCase.Execute(context.Background(), input)

		assert.NoError(t, err)
		_, err = uuid.Parse(output.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput.CustomerName, output.CustomerName)
		assert.Equal(t, expectedOutput.Status, output.Status)
		assert.Equal(t, expectedOutput.Total, output.Total)
		assert.Len(t, output.Items, len(expectedOutput.Items))

		mockRepo.AssertExpectations(t)
		mockPub.AssertExpectations(t)
	})

	t.Run("empty items", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		mockPub := new(usecasemock.MockOrderPublisher)
		useCase := usecase.NewCreateOrderUseCase(mockRepo, mockPub)

		input := dtos.OrderInput{
			CustomerName: "John Doe",
			Items:        []dtos.ItemInput{},
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, "order must contain at least one item", err.Error())
		assert.Empty(t, output)
	})

	t.Run("invalid item", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		mockPub := new(usecasemock.MockOrderPublisher)
		useCase := usecase.NewCreateOrderUseCase(mockRepo, mockPub)

		input := dtos.OrderInput{
			CustomerName: "John Doe",
			Items: []dtos.ItemInput{
				{
					ID:       "1",
					Name:     "Item 1",
					Quantity: 0,
					Price:    10.0,
				},
			},
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, "item quantity must be greater than zero", err.Error())
		assert.Empty(t, output)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(usecasemock.MockOrderRepository)
		mockPub := new(usecasemock.MockOrderPublisher)
		useCase := usecase.NewCreateOrderUseCase(mockRepo, mockPub)

		input := dtos.OrderInput{
			CustomerName: "John Doe",
			Items: []dtos.ItemInput{
				{
					ID:       "1",
					Name:     "Item 1",
					Quantity: 2,
					Price:    10.0,
				},
			},
		}

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(errors.New("database error"))

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Empty(t, output)
	})
}
