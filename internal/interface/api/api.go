package api

import (
	"order-service/internal/application/usecase"
)

type API struct {
	createOrderUseCase usecase.CreateOrderUseCase
	updateOrderUseCase usecase.UpdateOrderUseCase
	cancelOrderUseCase usecase.CancelOrderUseCase
	getOrderUseCase    usecase.GetOrderUseCase
	listOrderUseCase   usecase.ListOrderUseCase
}

func NewAPI(
	createOrderUseCase usecase.CreateOrderUseCase,
	updateOrderUseCase usecase.UpdateOrderUseCase,
	cancelOrderUseCase usecase.CancelOrderUseCase,
	getOrderUseCase usecase.GetOrderUseCase,
	listOrderUseCase usecase.ListOrderUseCase,
) *API {
	return &API{
		createOrderUseCase: createOrderUseCase,
		updateOrderUseCase: updateOrderUseCase,
		cancelOrderUseCase: cancelOrderUseCase,
		getOrderUseCase:    getOrderUseCase,
		listOrderUseCase:   listOrderUseCase,
	}
}
