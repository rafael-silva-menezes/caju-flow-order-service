package database_test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgresql://order_test_user:order_test_pass@localhost:5435/order_test_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			customer_name VARCHAR(255),
			status VARCHAR(15),
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS order_items (
			id VARCHAR(36) PRIMARY KEY,
			order_id VARCHAR(36) REFERENCES orders(id) ON DELETE CASCADE,
			name VARCHAR(255),
			quantity INT,
			price NUMERIC
		);
	`)
	require.NoError(t, err)

	_, err = db.Exec(`DELETE FROM order_items`)
	require.NoError(t, err)
	_, err = db.Exec(`DELETE FROM orders`)
	require.NoError(t, err)

	return db
}
