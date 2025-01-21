package database

import (
	"context"
	"database/sql"
	"errors"

	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"
)

type OrderRepositorySql struct {
	db *sql.DB
}

func NewOrderRepositorySql(db *sql.DB) *OrderRepositorySql {
	return &OrderRepositorySql{db: db}
}

func (r *OrderRepositorySql) Save(ctx context.Context, order *entity.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Inserir o pedido
	orderQuery := `
		INSERT INTO orders (id, customer_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET customer_name = $2, status = $3, updated_at = $5
	`
	_, err = tx.ExecContext(ctx, orderQuery, order.ID, order.CustomerName, order.Status.String(), order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return err
	}

	// Remover itens antigos
	itemDeleteQuery := `DELETE FROM order_items WHERE order_id = $1`
	_, err = tx.ExecContext(ctx, itemDeleteQuery, order.ID)
	if err != nil {
		return err
	}

	// Inserir itens novos
	itemInsertQuery := `
		INSERT INTO order_items (id, order_id, name, quantity, price)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemInsertQuery, item.ID, order.ID, item.Name, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepositorySql) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	// Buscar o pedido
	orderQuery := `
		SELECT id, customer_name, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, orderQuery, id)

	var order entity.Order
	var status string
	if err := row.Scan(&order.ID, &order.CustomerName, &status, &order.CreatedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	order.Status = parseOrderStatus(status)

	// Buscar itens do pedido
	itemQuery := `
		SELECT id, name, quantity, price
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := r.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *OrderRepositorySql) List(ctx context.Context) ([]entity.Order, error) {
	// Buscar todos os pedidos
	orderQuery := `
		SELECT id, customer_name, status, created_at, updated_at
		FROM orders
	`
	rows, err := r.db.QueryContext(ctx, orderQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var order entity.Order
		var status string
		if err := rows.Scan(&order.ID, &order.CustomerName, &status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		order.Status = parseOrderStatus(status)
		orderMap[order.ID] = &order
		orders = append(orders, order)
	}

	// Buscar itens para todos os pedidos
	itemQuery := `
		SELECT id, order_id, name, quantity, price
		FROM order_items
	`
	itemRows, err := r.db.QueryContext(ctx, itemQuery)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item entity.Item
		var orderID string
		if err := itemRows.Scan(&item.ID, &orderID, &item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		if order, exists := orderMap[orderID]; exists {
			order.Items = append(order.Items, item)
		}
	}

	return orders, nil
}

func parseOrderStatus(status string) entity.OrderStatus {
	switch status {
	case "pending":
		return entity.Pending
	case "processing":
		return entity.Processing
	case "completed":
		return entity.Completed
	case "canceled":
		return entity.Canceled
	default:
		return entity.Pending
	}
}
