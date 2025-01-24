package config

import (
	"database/sql"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func InitDatabase(cfg *Conf) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}

func InitQueue(cfg *Conf) (*amqp091.Connection, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.BrokerUser, cfg.BrokerPassword, cfg.BrokerHost, cfg.BrokerPort,
	)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to broker: %v", err)
	}
	return conn, nil
}
