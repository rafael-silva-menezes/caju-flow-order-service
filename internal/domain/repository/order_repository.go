package repository

import (
	"context"

	"order-service/internal/domain/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error

	FindByID(ctx context.Context, orderID string) (*entity.Order, error)

	List(ctx context.Context) ([]entity.Order, error)
}
