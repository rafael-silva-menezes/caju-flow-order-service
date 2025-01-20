package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/application/dtos"
	"order-service/internal/application/usecase"
	usecasemock "order-service/internal/application/usecase/mock"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateOrderUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		input       dtos.CreateOrderInput
		setupMocks  func(mockRepo *usecasemock.MockOrderRepository)
		expected    dtos.GetOrderOutput
		expectedErr error
	}{
		{
			name:    "should update pending order",
			orderID: "123",
			input: dtos.CreateOrderInput{
				CustomerName: "Jane",
				Items: []dtos.ItemInput{
					{ID: "item1", Name: "Item 1", Quantity: 2, Price: 10},
				},
			},
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				order := &entity.Order{
					OrderID:      "123",
					CustomerName: "John",
					Status:       entity.Pending,
					Items:        []entity.Item{{ID: "item1", Name: "Item 1", Quantity: 1, Price: 10}},
				}
				mockRepo.On("FindByID", mock.Anything, "123").Return(order, nil)
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			expected: dtos.GetOrderOutput{
				OrderID:      "123",
				CustomerName: "Jane",
				Total:        20,
				Status:       "pending",
				Items: []dtos.ItemOutput{
					{ID: "item1", Name: "Item 1", Quantity: 2, Price: 10, Total: 20},
				},
			},
			expectedErr: nil,
		},
		{
			name:    "should return error if order not found",
			orderID: "999",
			input: dtos.CreateOrderInput{
				CustomerName: "Jane",
				Items: []dtos.ItemInput{
					{ID: "item1", Name: "Item 1", Quantity: 1, Price: 10},
				},
			},
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				mockRepo.On("FindByID", mock.Anything, "999").Return(nil, repository.ErrNotFound)
			},
			expected:    dtos.GetOrderOutput{},
			expectedErr: errors.New("order not found"),
		},
		{
			name:    "should return error if order is not pending",
			orderID: "123",
			input: dtos.CreateOrderInput{
				CustomerName: "Jane",
				Items: []dtos.ItemInput{
					{ID: "item1", Name: "Item 1", Quantity: 1, Price: 10},
				},
			},
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				order := &entity.Order{
					OrderID:      "123",
					CustomerName: "John",
					Status:       entity.Completed,
					Items:        []entity.Item{{ID: "item1", Name: "Item 1", Quantity: 1, Price: 10}},
				}
				mockRepo.On("FindByID", mock.Anything, "123").Return(order, nil)
			},
			expected:    dtos.GetOrderOutput{},
			expectedErr: errors.New("order cannot be updated as it is not pending"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(usecasemock.MockOrderRepository)
			tt.setupMocks(mockRepo)
			updateOrderUseCase := usecase.NewUpdateOrderUseCase(mockRepo)

			result, err := updateOrderUseCase.Execute(context.Background(), tt.orderID, tt.input)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
