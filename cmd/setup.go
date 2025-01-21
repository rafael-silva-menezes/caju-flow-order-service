package main

import (
	"database/sql"
	"log"

	"order-service/internal/config"
	"order-service/internal/infrastructure/migrations"

	"github.com/streadway/amqp"
)

func setupInfra(cfg *config.Conf) (*sql.DB, *amqp.Connection, error) {
	db, err := config.InitDatabase(cfg)
	if err != nil {
		return nil, nil, err
	}

	if err := migrations.Migrate(db, "internal/infrastructure/migrations"); err != nil {
		return nil, nil, err
	}

	queueConn, err := config.InitQueue(cfg)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Infrastructure setup completed successfully")
	return db, queueConn, nil
}
