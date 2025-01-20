package dtos

import "order-service/internal/domain/entity"

type OrderInput struct {
	CustomerName string      `json:"customer_name"`
	Items        []ItemInput `json:"items"`
}

type OrderOutput struct {
	OrderID      string       `json:"order_id"`
	CustomerName string       `json:"customer_name"`
	Items        []ItemOutput `json:"items"`
	Total        float64      `json:"total"`
	Status       string       `json:"status"`
}

type ListOrderOutput struct {
	Orders []OrderOutput `json:"orders"`
}

func FromEntityToOrderOutput(order *entity.Order) OrderOutput {
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

	return OrderOutput{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		Total:        order.Total(),
		Status:       order.Status.String(),
		Items:        items,
	}
}
