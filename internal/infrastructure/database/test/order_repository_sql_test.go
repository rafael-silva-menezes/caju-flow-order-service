package database_test

import (
	"context"
	"testing"
	"time"

	"order-service/internal/domain/entity"
	"order-service/internal/domain/repository"
	"order-service/internal/infrastructure/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderRepositorySql_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewOrderRepositorySql(db)

	orderID := uuid.New().String()
	itemID1 := uuid.New().String()
	itemID2 := uuid.New().String()

	order := &entity.Order{
		ID:           orderID,
		CustomerName: "John Doe",
		Status:       entity.Pending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Items: []entity.Item{
			{ID: itemID1, Name: "Item 1", Quantity: 1, Price: 10.0},
			{ID: itemID2, Name: "Item 2", Quantity: 2, Price: 15.0},
		},
	}

	err := repo.Save(context.Background(), order)
	require.NoError(t, err)

	savedOrder, err := repo.FindByID(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, savedOrder.ID)
	assert.Equal(t, order.CustomerName, savedOrder.CustomerName)
	assert.Equal(t, order.Status, savedOrder.Status)
	assert.Equal(t, len(order.Items), len(savedOrder.Items))
}

func TestOrderRepositorySql_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewOrderRepositorySql(db)

	orderID := uuid.New().String()
	itemID := uuid.New().String()

	order := &entity.Order{
		ID:           orderID,
		CustomerName: "John Doe",
		Status:       entity.Pending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Items: []entity.Item{
			{ID: itemID, Name: "Item 1", Quantity: 1, Price: 10.0},
		},
	}

	err := repo.Save(context.Background(), order)
	require.NoError(t, err)

	savedOrder, err := repo.FindByID(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, savedOrder.ID)
	assert.Equal(t, order.CustomerName, savedOrder.CustomerName)
	assert.Equal(t, order.Status, savedOrder.Status)
	assert.Len(t, savedOrder.Items, 1)
}

func TestOrderRepositorySql_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewOrderRepositorySql(db)

	orderID1 := uuid.New().String()
	orderID2 := uuid.New().String()
	itemID1 := uuid.New().String()
	itemID2 := uuid.New().String()

	order1 := &entity.Order{
		ID:           orderID1,
		CustomerName: "John Doe",
		Status:       entity.Pending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Items: []entity.Item{
			{ID: itemID1, Name: "Item 1", Quantity: 1, Price: 10.0},
		},
	}
	order2 := &entity.Order{
		ID:           orderID2,
		CustomerName: "Jane Doe",
		Status:       entity.Completed,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Items: []entity.Item{
			{ID: itemID2, Name: "Item 2", Quantity: 2, Price: 15.0},
		},
	}

	err := repo.Save(context.Background(), order1)
	require.NoError(t, err)
	err = repo.Save(context.Background(), order2)
	require.NoError(t, err)

	orders, err := repo.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, orders, 2)
}

func TestOrderRepositorySql_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := database.NewOrderRepositorySql(db)

	order, err := repo.FindByID(context.Background(), "non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, repository.ErrNotFound, err)
}
