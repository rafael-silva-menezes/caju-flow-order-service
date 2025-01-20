package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/repository"
)

type GetOrderUseCase interface {
	Execute(ctx context.Context, orderID string) (dtos.GetOrderOutput, error)
}

type getOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewGetOrderUseCase(orderRepo repository.OrderRepository) GetOrderUseCase {
	return &getOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *getOrderUseCase) Execute(ctx context.Context, orderID string) (dtos.GetOrderOutput, error) {
	order, err := u.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dtos.GetOrderOutput{}, errors.New("order not found")
		}
		return dtos.GetOrderOutput{}, err
	}

	return dtos.FromEntityToGetOrderOutput(order), nil
}
