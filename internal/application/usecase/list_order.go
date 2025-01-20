package usecase

import (
	"context"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/repository"
)

type ListOrderUseCase interface {
	Execute(ctx context.Context) (dtos.ListOrderOutput, error)
}

type listOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewListOrderUseCase(orderRepo repository.OrderRepository) ListOrderUseCase {
	return &listOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *listOrderUseCase) Execute(ctx context.Context) (dtos.ListOrderOutput, error) {
	orders, err := u.orderRepository.List(ctx)
	if err != nil {
		return dtos.ListOrderOutput{}, err
	}

	var output dtos.ListOrderOutput
	for _, order := range orders {
		output.Orders = append(output.Orders, dtos.FromEntityToOrderOutput(&order))
	}

	return output, nil
}
