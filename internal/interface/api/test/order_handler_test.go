package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/internal/application/dtos"
	"order-service/internal/interface/api"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type mockCreateOrderUseCase struct {
	output dtos.OrderOutput
	err    error
}

func (m *mockCreateOrderUseCase) Execute(ctx context.Context, input dtos.OrderInput) (dtos.OrderOutput, error) {
	return m.output, m.err
}

func TestCreateOrder_Success(t *testing.T) {
	mockUseCase := &mockCreateOrderUseCase{
		output: dtos.OrderOutput{ID: "1", CustomerName: "John Doe", Status: "pending"},
	}
	api := api.NewAPI(mockUseCase, nil, nil, nil, nil)

	body := `{"customer_name": "John Doe", "items": [{"name": "item1", "quantity": 1, "price": 100}]}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()

	api.CreateOrder(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var response dtos.OrderOutput
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
}

func TestCreateOrder_InvalidInput(t *testing.T) {
	api := api.NewAPI(nil, nil, nil, nil, nil)

	body := `{"customer_name": "", "items": []}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()

	api.CreateOrder(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Customer name and items are required")
}

func TestCreateOrder_UseCaseError(t *testing.T) {
	mockUseCase := &mockCreateOrderUseCase{
		err: errors.New("use case error"),
	}
	api := api.NewAPI(mockUseCase, nil, nil, nil, nil)

	body := `{"customer_name": "John Doe", "items": [{"name": "item1", "quantity": 1, "price": 100}]}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()

	api.CreateOrder(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create order")
}

type mockGetOrderUseCase struct {
	output dtos.OrderOutput
	err    error
}

func (m *mockGetOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	return m.output, m.err
}

func TestGetOrder_Success(t *testing.T) {
	mockUseCase := &mockGetOrderUseCase{
		output: dtos.OrderOutput{ID: "1", CustomerName: "John Doe", Status: "pending"},
	}
	api := api.NewAPI(nil, nil, nil, mockUseCase, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/orders/{id}", api.GetOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response dtos.OrderOutput
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID)
}

func TestGetOrder_NotFound(t *testing.T) {
	mockUseCase := &mockGetOrderUseCase{
		err: errors.New("order not found"),
	}
	api := api.NewAPI(nil, nil, nil, mockUseCase, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/orders/{id}", api.GetOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Order not found")
}

type mockListOrderUseCase struct {
	output dtos.ListOrderOutput
	err    error
}

func (m *mockListOrderUseCase) Execute(ctx context.Context) (dtos.ListOrderOutput, error) {
	return m.output, m.err
}

func TestListOrders_Success(t *testing.T) {
	mockUseCase := &mockListOrderUseCase{
		output: dtos.ListOrderOutput{
			Orders: []dtos.OrderOutput{
				{ID: "1", CustomerName: "John Doe", Status: "pending"},
				{ID: "2", CustomerName: "Jane Doe", Status: "completed"},
			},
		},
	}
	api := api.NewAPI(nil, nil, nil, nil, mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/orders", api.ListOrders)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response dtos.ListOrderOutput
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Orders, 2)
	assert.Equal(t, "1", response.Orders[0].ID)
}

func TestListOrders_Failure(t *testing.T) {
	mockUseCase := &mockListOrderUseCase{
		err: errors.New("failed to fetch orders"),
	}
	api := api.NewAPI(nil, nil, nil, nil, mockUseCase)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/orders", api.ListOrders)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to list orders")
}

type mockUpdateOrderUseCase struct {
	output dtos.OrderOutput
	err    error
}

func (m *mockUpdateOrderUseCase) Execute(ctx context.Context, id string, input dtos.OrderInput) (dtos.OrderOutput, error) {
	return m.output, m.err
}

func TestUpdateOrder_Success(t *testing.T) {
	mockUseCase := &mockUpdateOrderUseCase{
		output: dtos.OrderOutput{ID: "1", CustomerName: "John Doe", Status: "completed"},
	}
	api := api.NewAPI(nil, mockUseCase, nil, nil, nil)

	body := `{"customer_name": "John Doe", "items": [{"name": "item1", "quantity": 2, "price": 200}]}`
	req := httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/orders/{id}", api.UpdateOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response dtos.OrderOutput
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "completed", response.Status)
}

func TestUpdateOrder_InvalidInput(t *testing.T) {
	mockUseCase := &mockUpdateOrderUseCase{
		err: errors.New("use case error"),
	}

	api := api.NewAPI(nil, mockUseCase, nil, nil, nil)

	body := `{"customer_name": "", "items": []}`
	req := httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewReader([]byte(body)))
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/orders/{id}", api.UpdateOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Customer name and items are required")
}

type mockCancelOrderUseCase struct {
	output dtos.OrderOutput
	err    error
}

func (m *mockCancelOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	return m.output, m.err
}

func TestCancelOrder_Success(t *testing.T) {
	mockUseCase := &mockCancelOrderUseCase{
		output: dtos.OrderOutput{ID: "1", CustomerName: "John Doe", Status: "canceled"},
	}
	api := api.NewAPI(nil, nil, mockUseCase, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/orders/{id}", api.CancelOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response dtos.OrderOutput
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "canceled", response.Status)
}

func TestCancelOrder_Failure(t *testing.T) {
	mockUseCase := &mockCancelOrderUseCase{
		err: errors.New("failed to cancel order"),
	}
	api := api.NewAPI(nil, nil, mockUseCase, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/orders/{id}", api.CancelOrder)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to cancel order")
}
