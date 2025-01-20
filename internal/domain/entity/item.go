package domain

import (
	"errors"
	"fmt"
)

type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func NewItem(id, name string, quantity int, price float64) (*Item, error) {
	if name == "" {
		return nil, errors.New("item name cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("item quantity must be greater than zero")
	}
	if price <= 0 {
		return nil, errors.New("item price must be greater than zero")
	}

	return &Item{
		ID:       id,
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}, nil
}

func (i *Item) Total() float64 {
	return float64(i.Quantity) * i.Price
}

func (i *Item) UpdateQuantity(newQuantity int) error {
	if newQuantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	i.Quantity = newQuantity
	return nil
}

func (i *Item) UpdatePrice(newPrice float64) error {
	if newPrice <= 0 {
		return errors.New("price must be greater than zero")
	}
	i.Price = newPrice
	return nil
}

func (i *Item) String() string {
	return fmt.Sprintf("Item{id: %s, name: %s, quantity: %d, price: %.2f}", i.ID, i.Name, i.Quantity, i.Price)
}
