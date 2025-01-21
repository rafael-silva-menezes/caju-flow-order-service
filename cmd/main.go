package main

import (
	"log"

	"order-service/internal/config"
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

	log.Println("Application started...")
}
