package entity

import (
	"errors"
	"time"
)

type OrderStatus int

const (
	Pending OrderStatus = iota
	Processing
	Completed
	Canceled
)

func (s OrderStatus) String() string {
	return [...]string{"pending", "processing", "completed", "canceled"}[s]
}

type Order struct {
	OrderID      string      `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Items        []Item      `json:"items"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func NewOrder(orderID, customerName string, items []Item) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	for _, item := range items {
		if item.Quantity <= 0 || item.Price <= 0 {
			return nil, errors.New("invalid item quantity or price")
		}
	}

	return &Order{
		OrderID:      orderID,
		CustomerName: customerName,
		Items:        items,
		Status:       Pending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (o *Order) IsValid() error {
	if o.OrderID == "" {
		return errors.New("invalid order id")
	}
	return nil
}

func (o *Order) Total() float64 {
	var total float64
	for _, item := range o.Items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}

func (o *Order) UpdateOrderDetails(newCustomerName string, newItems []Item) error {
	if o.Status != Pending {
		return errors.New("order cannot be modified as it is not pending")
	}

	if newCustomerName != "" {
		o.CustomerName = newCustomerName
	}

	if len(newItems) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range newItems {
		if item.Quantity <= 0 || item.Price <= 0 {
			return errors.New("invalid item quantity or price")
		}
	}

	o.Items = newItems
	o.UpdatedAt = time.Now()

	return nil
}

// Método SetStatus para atualizar o status do pedido
func (o *Order) SetStatus(status OrderStatus) error {
	// Verifica se o status atual permite a alteração
	if o.Status == status {
		return nil // Caso o status seja o mesmo, não faz sentido alterá-lo
	}

	// Caso o pedido já tenha sido processado ou concluído, ele não pode ser alterado para outro status
	if o.Status == Completed || o.Status == Canceled {
		return errors.New("cannot change status of completed or canceled order")
	}

	o.Status = status
	o.UpdatedAt = time.Now()
	return nil
}
