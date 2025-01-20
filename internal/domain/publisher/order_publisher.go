package publisher

import (
	"context"

	"order-service/internal/domain/entity"
)

type OrderPublisher interface {
	PublishCreatedOrder(ctx context.Context, order *entity.Order) error
}
