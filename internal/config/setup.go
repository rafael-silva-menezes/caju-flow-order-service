package config

import (
	"database/sql"
	"log"

	"order-service/internal/infrastructure/migrations"

	"github.com/streadway/amqp"
)

func SetupInfra(cfg *Conf) (*sql.DB, *amqp.Connection, error) {
	db, err := InitDatabase(cfg)
	if err != nil {
		return nil, nil, err
	}

	if err := migrations.Migrate(db, "internal/infrastructure/migrations"); err != nil {
		return nil, nil, err
	}

	queueConn, err := InitQueue(cfg)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Infrastructure setup completed successfully")
	return db, queueConn, nil
}
