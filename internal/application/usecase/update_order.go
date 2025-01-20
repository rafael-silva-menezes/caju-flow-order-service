package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"
)

type UpdateOrderUseCase interface {
	Execute(ctx context.Context, orderID string, input dtos.CreateOrderInput) (dtos.GetOrderOutput, error)
}

type updateOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewUpdateOrderUseCase(orderRepo repository.OrderRepository) UpdateOrderUseCase {
	return &updateOrderUseCase{
		orderRepository: orderRepo,
	}
}

func (u *updateOrderUseCase) Execute(ctx context.Context, orderID string, input dtos.CreateOrderInput) (dtos.GetOrderOutput, error) {
	if input.CustomerName == "" {
		return dtos.GetOrderOutput{}, errors.New("customer name cannot be empty")
	}
	if len(input.Items) == 0 {
		return dtos.GetOrderOutput{}, errors.New("order must contain at least one item")
	}

	order, err := u.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return dtos.GetOrderOutput{}, errors.New("order not found")
		}
		return dtos.GetOrderOutput{}, err
	}

	if order.Status != entity.Pending {
		return dtos.GetOrderOutput{}, errors.New("order cannot be updated as it is not pending")
	}

	var items []entity.Item
	for _, itemInput := range input.Items {
		item, err := entity.NewItem(itemInput.ID, itemInput.Name, itemInput.Quantity, itemInput.Price)
		if err != nil {
			return dtos.GetOrderOutput{}, err
		}
		items = append(items, *item)
	}

	err = order.UpdateOrderDetails(input.CustomerName, items)
	if err != nil {
		return dtos.GetOrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, order)
	if err != nil {
		return dtos.GetOrderOutput{}, err
	}

	return dtos.FromEntityToGetOrderOutput(order), nil
}
