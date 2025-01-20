package order

import (
	"context"
)

type GetOrderUseCase interface {
	Execute(ctx context.Context, orderID string) (GetOrderOutput, error)
}

type GetOrderOutput struct {
	OrderID      string  `json:"order_id"`
	CustomerName string  `json:"customer_name"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
}
