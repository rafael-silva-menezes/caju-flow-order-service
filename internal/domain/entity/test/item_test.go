package entity_test

import (
	"testing"

	"order-service/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItem_ValidItem(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)

	require.NoError(t, err)

	assert.Equal(t, "1", item.ID)
	assert.Equal(t, "Product A", item.Name)
	assert.Equal(t, 2, item.Quantity)
	assert.Equal(t, 50.0, item.Price)
}

func TestNewItem_InvalidItem(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", -2, 50.0)

	require.Error(t, err)
	assert.Nil(t, item)
}

func TestNewItem_EmptyName(t *testing.T) {
	item, err := entity.NewItem("1", "", 2, 50.0)

	require.Error(t, err)
	assert.Nil(t, item)
}

func TestItem_Total(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)
	require.NoError(t, err)

	assert.Equal(t, 100.0, item.Total())
}

func TestItem_UpdateQuantity(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)
	require.NoError(t, err)

	err = item.UpdateQuantity(3)
	require.NoError(t, err)

	assert.Equal(t, 3, item.Quantity)
}

func TestItem_UpdateQuantity_Invalid(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)
	require.NoError(t, err)

	err = item.UpdateQuantity(-1)
	require.Error(t, err)
}

func TestItem_UpdatePrice(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)
	require.NoError(t, err)

	err = item.UpdatePrice(60.0)
	require.NoError(t, err)

	assert.Equal(t, 60.0, item.Price)
}

func TestItem_UpdatePrice_Invalid(t *testing.T) {
	item, err := entity.NewItem("1", "Product A", 2, 50.0)
	require.NoError(t, err)

	err = item.UpdatePrice(-10.0)
	require.Error(t, err)
}
