package main

import (
	"log"
	"net/http"

	"order-service/internal/application/usecase"
	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/publisher"
	"order-service/internal/interface/api"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	db, queueConn, err := config.SetupInfra(cfg)
	if err != nil {
		log.Fatalf("error setting up infrastructure: %v", err)
	}
	defer db.Close()
	defer queueConn.Close()

	orderRepository := database.NewOrderRepositorySql(db)
	rabbitPublisher, err := publisher.NewRabbitMQPublisher(queueConn, "orders_exchange", "order_created", "order_created_queue")
	if err != nil {
		log.Fatalf("error creating RabbitMQ publisher: %v", err)
	}

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, rabbitPublisher)
	updateOrderUseCase := usecase.NewUpdateOrderUseCase(orderRepository)
	cancelOrderUseCase := usecase.NewCancelOrderUseCase(orderRepository)
	getOrderUseCase := usecase.NewGetOrderUseCase(orderRepository)
	listOrderUseCase := usecase.NewListOrderUseCase(orderRepository)

	handlers := api.NewAPI(
		createOrderUseCase,
		updateOrderUseCase,
		cancelOrderUseCase,
		getOrderUseCase,
		listOrderUseCase,
	)

	r := api.NewRouter(handlers)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
