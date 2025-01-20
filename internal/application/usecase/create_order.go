package usecase

import (
	"context"
	"errors"

	"order-service/internal/application/dtos"
	"order-service/internal/domain/entity"
	"order-service/internal/domain/publisher"
	"order-service/internal/domain/repository"

	"github.com/google/uuid"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, input dtos.OrderInput) (dtos.OrderOutput, error)
}

type createOrderUseCase struct {
	orderRepository repository.OrderRepository
	orderPublisher  publisher.OrderPublisher
}

func NewCreateOrderUseCase(orderRepo repository.OrderRepository, orderPub publisher.OrderPublisher) CreateOrderUseCase {
	return &createOrderUseCase{
		orderRepository: orderRepo,
		orderPublisher:  orderPub,
	}
}

func (u *createOrderUseCase) Execute(ctx context.Context, input dtos.OrderInput) (dtos.OrderOutput, error) {
	if len(input.Items) == 0 {
		return dtos.OrderOutput{}, errors.New("order must contain at least one item")
	}

	var items []entity.Item
	for _, itemInput := range input.Items {
		item, err := entity.NewItem(itemInput.ID, itemInput.Name, itemInput.Quantity, itemInput.Price)
		if err != nil {
			return dtos.OrderOutput{}, err
		}
		items = append(items, *item)
	}

	newOrder, err := entity.NewOrder(generateOrderID(), input.CustomerName, items)
	if err != nil {
		return dtos.OrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, newOrder)
	if err != nil {
		return dtos.OrderOutput{}, err
	}

	err = u.orderPublisher.PublishCreatedOrder(ctx, newOrder)
	if err != nil {
		return dtos.OrderOutput{}, err
	}

	return dtos.FromEntityToOrderOutput(newOrder), nil
}

func generateOrderID() string {
	return uuid.New().String()
}
