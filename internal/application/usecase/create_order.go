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
	Execute(ctx context.Context, input dtos.CreateOrderInput) (dtos.CreateOrderOutput, error)
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

func (u *createOrderUseCase) Execute(ctx context.Context, input dtos.CreateOrderInput) (dtos.CreateOrderOutput, error) {
	if len(input.Items) == 0 {
		return dtos.CreateOrderOutput{}, errors.New("order must contain at least one item")
	}

	var items []entity.Item
	for _, itemInput := range input.Items {
		item, err := entity.NewItem(itemInput.ID, itemInput.Name, itemInput.Quantity, itemInput.Price)
		if err != nil {
			return dtos.CreateOrderOutput{}, err
		}
		items = append(items, *item)
	}

	newOrder, err := entity.NewOrder(generateOrderID(), input.CustomerName, items)
	if err != nil {
		return dtos.CreateOrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, newOrder)
	if err != nil {
		return dtos.CreateOrderOutput{}, err
	}

	err = u.orderPublisher.PublishCreatedOrder(ctx, newOrder)
	if err != nil {
		return dtos.CreateOrderOutput{}, err
	}

	output := dtos.CreateOrderOutput{
		OrderID:      newOrder.OrderID,
		CustomerName: newOrder.CustomerName,
		Total:        newOrder.Total(),
		Status:       newOrder.Status.String(),
	}
	for _, item := range newOrder.Items {
		output.Items = append(output.Items, dtos.ItemOutput{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			Total:    item.Total(),
		})
	}

	return output, nil
}

func generateOrderID() string {
	return uuid.New().String()
}
