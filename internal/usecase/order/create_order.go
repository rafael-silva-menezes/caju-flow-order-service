package order

import (
	"context"
)

type CreateOrderInput struct {
	CustomerName string      `json:"customer_name"`
	Items        []ItemInput `json:"items"`
}

type ItemInput struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type CreateOrderOutput struct {
	OrderID      string       `json:"order_id"`
	CustomerName string       `json:"customer_name"`
	Items        []ItemOutput `json:"items"`
	Total        float64      `json:"total"`
	Status       string       `json:"status"`
}

type ItemOutput struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total"`
}

type CreateOrderUseCase interface {
	Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error)
}
