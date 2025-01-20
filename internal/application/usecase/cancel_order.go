package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"
)

type CancelOrderUseCase interface {
	Execute(ctx context.Context, orderID string) (dtos.GetOrderOutput, error)
}

type cancelOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewCancelOrderUseCase(orderRepo repository.OrderRepository) CancelOrderUseCase {
	return &cancelOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *cancelOrderUseCase) Execute(ctx context.Context, orderID string) (dtos.GetOrderOutput, error) {
	order, err := u.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dtos.GetOrderOutput{}, errors.New("order not found")
		}
		return dtos.GetOrderOutput{}, err
	}

	if order.Status != entity.Pending {
		return dtos.GetOrderOutput{}, errors.New("only pending orders can be canceled")
	}

	err = order.SetStatus(entity.Canceled)
	if err != nil {
		return dtos.GetOrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, order)
	if err != nil {
		return dtos.GetOrderOutput{}, err
	}

	return dtos.FromEntityToGetOrderOutput(order), nil
}
