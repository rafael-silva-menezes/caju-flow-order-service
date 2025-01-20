package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/repository"
)

type GetOrderUseCase interface {
	Execute(ctx context.Context, id string) (dtos.OrderOutput, error)
}

type getOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewGetOrderUseCase(orderRepo repository.OrderRepository) GetOrderUseCase {
	return &getOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *getOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	order, err := u.orderRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dtos.OrderOutput{}, errors.New("order not found")
		}
		return dtos.OrderOutput{}, err
	}

	return dtos.FromEntityToOrderOutput(order), nil
}
