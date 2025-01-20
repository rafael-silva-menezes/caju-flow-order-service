package order_test

import (
	"context"

	"order-service/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type mockOrderRepository struct {
	mock.Mock
}

func (m *mockOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockOrderRepository) List(ctx context.Context) ([]entity.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Order), args.Error(1)
}

type mockOrderPublisher struct {
	mock.Mock
}

func (m *mockOrderPublisher) PublishCreatedOrder(ctx context.Context, order *entity.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
