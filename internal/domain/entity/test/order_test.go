package entity_test

import (
	"testing"

	"order-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrder_ValidOrder(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
		{ID: "2", Name: "Product B", Quantity: 1, Price: 30.0},
	}

	order, err := entity.NewOrder("12345", "João Silva", items)

	require.NoError(t, err)

	assert.Equal(t, "12345", order.OrderID)
	assert.Equal(t, "João Silva", order.CustomerName)
	assert.Equal(t, entity.Pending, order.Status)
	assert.Len(t, order.Items, 2)
	assert.Equal(t, 130.0, order.Total())
}

func TestNewOrder_InvalidItems(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: -2, Price: 50.0},
	}

	order, err := entity.NewOrder("12345", "João Silva", items)

	require.Error(t, err)
	assert.Nil(t, order)
}

func TestNewOrder_EmptyItems(t *testing.T) {
	order, err := entity.NewOrder("12345", "João Silva", nil)

	require.Error(t, err)
	assert.Nil(t, order)
}

func TestOrder_IsValid(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	err = order.IsValid()
	assert.NoError(t, err)

	order.OrderID = ""
	err = order.IsValid()
	assert.Error(t, err)
}

func TestOrder_Total(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
		{ID: "2", Name: "Product B", Quantity: 1, Price: 30.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	assert.Equal(t, 130.0, order.Total())
}

func TestOrder_UpdateOrderDetails(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
		{ID: "2", Name: "Product B", Quantity: 1, Price: 30.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	newItems := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 3, Price: 45.0},
		{ID: "3", Name: "Product C", Quantity: 2, Price: 25.0},
	}

	err = order.UpdateOrderDetails("Maria Silva", newItems)
	require.NoError(t, err)

	assert.Equal(t, "Maria Silva", order.CustomerName)
	assert.Len(t, order.Items, 2)
	assert.Equal(t, 3, order.Items[0].Quantity)
	assert.Equal(t, 2, order.Items[1].Quantity)
	assert.Equal(t, "Product C", order.Items[1].Name)

	assert.True(t, order.UpdatedAt.After(order.CreatedAt))
}

func TestOrder_UpdateOrderDetails_NotPending(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	order.Status = entity.Processing

	newItems := []entity.Item{
		{ID: "2", Name: "Product B", Quantity: 2, Price: 40.0},
	}

	err = order.UpdateOrderDetails("Maria Silva", newItems)
	require.Error(t, err)
	assert.Equal(t, "order cannot be modified as it is not pending", err.Error())
}

func TestOrder_UpdateOrderDetails_InvalidItem(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	newItems := []entity.Item{
		{ID: "2", Name: "Product B", Quantity: -1, Price: 40.0},
	}

	err = order.UpdateOrderDetails("Maria Silva", newItems)
	require.Error(t, err)
	assert.Equal(t, "invalid item quantity or price", err.Error())
}

func TestOrder_UpdateOrderDetails_EmptyItems(t *testing.T) {
	items := []entity.Item{
		{ID: "1", Name: "Product A", Quantity: 2, Price: 50.0},
	}
	order, err := entity.NewOrder("12345", "João Silva", items)
	require.NoError(t, err)

	err = order.UpdateOrderDetails("Maria Silva", nil)
	require.Error(t, err)
	assert.Equal(t, "order must contain at least one item", err.Error()) // Atualizado
}
