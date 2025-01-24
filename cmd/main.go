package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"order-service/internal/application/usecase"
	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/publisher"
	"order-service/internal/interface/api"
)

func main() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "" {
		env = "prod"
	}

	log.Printf("Running in %s environment\n", env)

	var logger *log.Logger
	if env == "dev" {
		logger = log.New(os.Stdout, "DEV ", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stdout, "PROD ", log.LstdFlags)
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		logger.Fatalf("Error loading configuration: %v", err)
	}

	db, queueConn, err := config.SetupInfra(cfg)
	if err != nil {
		logger.Fatalf("Error setting up infrastructure: %v", err)
	}
	defer db.Close()
	defer queueConn.Close()

	orderRepository := database.NewOrderRepositorySql(db)
	rabbitPublisher, err := publisher.NewRabbitMQPublisher(queueConn, "orders_exchange", "order_created", "order_created_queue")
	if err != nil {
		logger.Fatalf("Error creating RabbitMQ publisher: %v", err)
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

	logger.Println("Starting server on :8080...")
	logger.Fatal(http.ListenAndServe(":8080", r))
}
