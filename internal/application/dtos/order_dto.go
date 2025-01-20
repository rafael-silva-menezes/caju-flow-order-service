package dtos

import "order-service/internal/domain/entity"

type CreateOrderInput struct {
	CustomerName string      `json:"customer_name"`
	Items        []ItemInput `json:"items"`
}

type CreateOrderOutput struct {
	OrderID      string       `json:"order_id"`
	CustomerName string       `json:"customer_name"`
	Items        []ItemOutput `json:"items"`
	Total        float64      `json:"total"`
	Status       string       `json:"status"`
}

type GetOrderOutput struct {
	OrderID      string       `json:"order_id"`
	CustomerName string       `json:"customer_name"`
	Total        float64      `json:"total"`
	Status       string       `json:"status"`
	Items        []ItemOutput `json:"items"`
}

type OrderOutput struct {
	OrderID      string  `json:"order_id"`
	CustomerName string  `json:"customer_name"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
}

type ListOrderOutput struct {
	Orders []OrderOutput `json:"orders"`
}

func FromEntityToGetOrderOutput(order *entity.Order) GetOrderOutput {
	var items []ItemOutput
	for _, item := range order.Items {
		items = append(items, ItemOutput{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			Total:    item.Total(),
		})
	}

	return GetOrderOutput{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Total:        order.Total(),
		Status:       order.Status.String(),
		Items:        items,
	}
}

func FromEntityToOrderOutput(order *entity.Order) OrderOutput {
	return OrderOutput{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Total:        order.Total(),
		Status:       order.Status.String(),
	}
}

func FromEntityToCreateOrderOutput(order *entity.Order) CreateOrderOutput {
	var items []ItemOutput
	for _, item := range order.Items {
		items = append(items, ItemOutput{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			Total:    item.Total(),
		})
	}

	return CreateOrderOutput{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Total:        order.Total(),
		Status:       order.Status.String(),
		Items:        items,
	}
}
