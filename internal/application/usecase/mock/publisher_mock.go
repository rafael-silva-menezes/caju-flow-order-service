package usecase_mock

import (
	"context"

	"order-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockOrderPublisher struct {
	mock.Mock
}

func (m *MockOrderPublisher) PublishCreatedOrder(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
