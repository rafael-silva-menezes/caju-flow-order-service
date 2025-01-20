package dtos

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
	OrderID      string  `json:"order_id"`
	CustomerName string  `json:"customer_name"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
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
