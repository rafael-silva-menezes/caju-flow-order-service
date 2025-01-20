package order

import (
	"context"
	"errors"

	"order-service/internal/domain/repository"
)

type getOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewGetOrderUseCase(orderRepo repository.OrderRepository) GetOrderUseCase {
	return &getOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *getOrderUseCase) Execute(ctx context.Context, orderID string) (GetOrderOutput, error) {
	order, err := u.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return GetOrderOutput{}, errors.New("order not found")
		}
		return GetOrderOutput{}, err
	}

	return GetOrderOutput{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Total:        order.Total(),
		Status:       order.Status.String(),
	}, nil
}
