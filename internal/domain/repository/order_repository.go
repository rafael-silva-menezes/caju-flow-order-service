package repository

import (
	"context"
	"errors"

	"order-service/internal/domain/entity"
)

var ErrNotFound = errors.New("order not found")

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error

	FindByID(ctx context.Context, orderID string) (*entity.Order, error)

	List(ctx context.Context) ([]entity.Order, error)
}
