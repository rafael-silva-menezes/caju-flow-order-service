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

func TestCancelOrderUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMocks  func(mockRepo *usecasemock.MockOrderRepository)
		expected    dtos.OrderOutput
		expectedErr error
	}{
		{
			name: "should cancel pending order",
			id:   "123",
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				order := &entity.Order{
					ID:           "123",
					CustomerName: "John",
					Status:       entity.Pending,
				}
				mockRepo.On("FindByID", mock.Anything, "123").Return(order, nil)
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			expected: dtos.OrderOutput{
				ID:           "123",
				CustomerName: "John",
				Total:        0,
				Status:       "canceled",
			},
			expectedErr: nil,
		},
		{
			name: "should return error if order not found",
			id:   "999",
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				mockRepo.On("FindByID", mock.Anything, "999").Return(nil, repository.ErrNotFound)
			},
			expected:    dtos.OrderOutput{},
			expectedErr: errors.New("order not found"),
		},
		{
			name: "should return error if order is not pending",
			id:   "123",
			setupMocks: func(mockRepo *usecasemock.MockOrderRepository) {
				order := &entity.Order{
					ID:           "123",
					CustomerName: "John",
					Status:       entity.Completed,
				}
				mockRepo.On("FindByID", mock.Anything, "123").Return(order, nil)
			},
			expected:    dtos.OrderOutput{},
			expectedErr: errors.New("only pending orders can be canceled"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(usecasemock.MockOrderRepository)
			tt.setupMocks(mockRepo)
			cancelOrderUseCase := usecase.NewCancelOrderUseCase(mockRepo)

			result, err := cancelOrderUseCase.Execute(context.Background(), tt.id)

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
