package order

import (
	"context"
	"errors"
	"time"

	"order-service/internal/domain/entity"
	"order-service/internal/domain/publisher"
	"order-service/internal/domain/repository"
)

type createOrderUseCase struct {
	orderRepository repository.OrderRepository
	orderPublisher  publisher.OrderPublisher
}

// NewCreateOrderUseCase creates a new instance of the CreateOrderUseCase.
func NewCreateOrderUseCase(orderRepo repository.OrderRepository, orderPub publisher.OrderPublisher) CreateOrderUseCase {
	return &createOrderUseCase{
		orderRepository: orderRepo,
		orderPublisher:  orderPub,
	}
}

// Execute handles the logic for creating a new order.
func (u *createOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	if len(input.Items) == 0 {
		return CreateOrderOutput{}, errors.New("order must contain at least one item")
	}

	var items []entity.Item
	for _, itemInput := range input.Items {
		item, err := entity.NewItem(itemInput.ID, itemInput.Name, itemInput.Quantity, itemInput.Price)
		if err != nil {
			return CreateOrderOutput{}, err
		}
		items = append(items, *item)
	}

	newOrder, err := entity.NewOrder(generateOrderID(), input.CustomerName, items)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	err = u.orderRepository.Save(ctx, newOrder)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	err = u.orderPublisher.PublishCreatedOrder(ctx, newOrder)
	if err != nil {
		return CreateOrderOutput{}, err
	}

	output := CreateOrderOutput{
		OrderID:      newOrder.OrderID,
		CustomerName: newOrder.CustomerName,
		Total:        newOrder.Total(),
		Status:       newOrder.Status.String(),
	}
	for _, item := range newOrder.Items {
		output.Items = append(output.Items, ItemOutput{
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
	return time.Now().Format("20060102150405")
}
