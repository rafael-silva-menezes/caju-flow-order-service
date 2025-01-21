package database_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"order-service/internal/domain/entity"
	"order-service/internal/infrastructure/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Load environment variables from .env.test file
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	// Set up database connection for testing
	dbURL := os.Getenv("DB_URL")
	println(dbURL)
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	testDB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	defer testDB.Close()

	// Create test tables
	if err := setupDatabase(); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	// Run tests
	code := m.Run()

	// Clean up test tables
	teardownDatabase()

	os.Exit(code)
}

func setupDatabase() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			customer_name TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			price NUMERIC(10,2) NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
		)`,
	}
	for _, query := range queries {
		if _, err := testDB.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func teardownDatabase() {
	queries := []string{
		`DROP TABLE IF EXISTS order_items`,
		`DROP TABLE IF EXISTS orders`,
	}
	for _, query := range queries {
		_, _ = testDB.Exec(query) // Ignorar erros ao limpar
	}
}

func TestOrderRepository_Save(t *testing.T) {
	repo := database.NewOrderRepositorySql(testDB)
	ctx := context.Background()

	orderItems := []entity.Item{
		{ID: "item1", Name: "Item 1", Quantity: 2, Price: 50.0},
		{ID: "item2", Name: "Item 2", Quantity: 1, Price: 100.0},
	}
	order, err := entity.NewOrder("order1", "John Doe", orderItems)
	assert.NoError(t, err)

	// Salvar o pedido
	err = repo.Save(ctx, order)
	assert.NoError(t, err)

	// Validar se o pedido foi salvo corretamente
	savedOrder, err := repo.FindByID(ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order.ID, savedOrder.ID)
	assert.Equal(t, order.CustomerName, savedOrder.CustomerName)
	assert.Equal(t, order.Status, savedOrder.Status)
	assert.Equal(t, len(order.Items), len(savedOrder.Items))
}

func TestOrderRepository_FindByID(t *testing.T) {
	repo := database.NewOrderRepositorySql(testDB)
	ctx := context.Background()

	// Criar pedido para buscar
	orderItems := []entity.Item{
		{ID: "item1", Name: "Item 1", Quantity: 2, Price: 50.0},
	}
	order, err := entity.NewOrder("order2", "Jane Doe", orderItems)
	assert.NoError(t, err)

	err = repo.Save(ctx, order)
	assert.NoError(t, err)

	// Buscar o pedido
	foundOrder, err := repo.FindByID(ctx, order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order.ID, foundOrder.ID)
	assert.Equal(t, order.CustomerName, foundOrder.CustomerName)
	assert.Equal(t, order.Status, foundOrder.Status)
}

func TestOrderRepository_List(t *testing.T) {
	repo := database.NewOrderRepositorySql(testDB)
	ctx := context.Background()

	// Criar pedidos
	order1, _ := entity.NewOrder("order3", "Alice", []entity.Item{
		{ID: "item1", Name: "Item 1", Quantity: 1, Price: 10.0},
	})
	order2, _ := entity.NewOrder("order4", "Bob", []entity.Item{
		{ID: "item2", Name: "Item 2", Quantity: 2, Price: 20.0},
	})

	err := repo.Save(ctx, order1)
	assert.NoError(t, err)
	err = repo.Save(ctx, order2)
	assert.NoError(t, err)

	// Listar pedidos
	orders, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
}
