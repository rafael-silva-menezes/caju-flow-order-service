package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"
)

type CancelOrderUseCase interface {
	Execute(ctx context.Context, id string) (dtos.OrderOutput, error)
}

type cancelOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewCancelOrderUseCase(orderRepo repository.OrderRepository) CancelOrderUseCase {
	return &cancelOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *cancelOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	order, err := u.orderRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dtos.OrderOutput{}, errors.New("order not found")
		}
		return dtos.OrderOutput{}, err
	}

	if order.Status != entity.Pending {
		return dtos.OrderOutput{}, errors.New("only pending orders can be canceled")
	}

	err = order.SetStatus(entity.Canceled)
	if err != nil {
		return dtos.OrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, order)
	if err != nil {
		return dtos.OrderOutput{}, err
	}

	return dtos.FromEntityToOrderOutput(order), nil
}
