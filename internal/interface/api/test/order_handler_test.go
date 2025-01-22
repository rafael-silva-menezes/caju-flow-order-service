package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/internal/application/dtos"
	"order-service/internal/application/usecase"
	"order-service/internal/interface/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCreateOrderUseCase struct {
	mock.Mock
}

func (m *MockCreateOrderUseCase) Execute(ctx context.Context, input dtos.OrderInput) (dtos.OrderOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dtos.OrderOutput), args.Error(1)
}

type MockGetOrderUseCase struct {
	mock.Mock
}

func (m *MockGetOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.OrderOutput), args.Error(1)
}

type MockListOrderUseCase struct {
	mock.Mock
}

func (m *MockListOrderUseCase) Execute(ctx context.Context) (dtos.ListOrderOutput, error) {
	args := m.Called(ctx)
	return args.Get(0).(dtos.ListOrderOutput), args.Error(1)
}

type MockUpdateOrderUseCase struct {
	mock.Mock
}

func (m *MockUpdateOrderUseCase) Execute(ctx context.Context, id string, input dtos.OrderInput) (dtos.OrderOutput, error) {
	args := m.Called(ctx, id, input)
	return args.Get(0).(dtos.OrderOutput), args.Error(1)
}

type MockCancelOrderUseCase struct {
	mock.Mock
}

func (m *MockCancelOrderUseCase) Execute(ctx context.Context, id string) (dtos.OrderOutput, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.OrderOutput), args.Error(1)
}

func setupRouter(
	createOrderUseCase usecase.CreateOrderUseCase,
	updateOrderUseCase usecase.UpdateOrderUseCase,
	cancelOrderUseCase usecase.CancelOrderUseCase,
	getOrderUseCase usecase.GetOrderUseCase,
	listOrderUseCase usecase.ListOrderUseCase,
) *api.API {
	return api.NewAPI(
		createOrderUseCase,
		updateOrderUseCase,
		cancelOrderUseCase,
		getOrderUseCase,
		listOrderUseCase,
	)
}

func TestCreateOrder(t *testing.T) {
	mockCreateOrderUseCase := new(MockCreateOrderUseCase)
	expectedOrder := dtos.OrderOutput{
		ID: "1",
	}

	mockCreateOrderUseCase.On("Execute", mock.Anything, mock.Anything).Return(expectedOrder, nil)

	apiInstance := setupRouter(mockCreateOrderUseCase, nil, nil, nil, nil)

	orderInput := dtos.OrderInput{
		CustomerName: "John Doe",
		Items: []dtos.ItemInput{
			{
				ID:       "1",
				Name:     "Item 1",
				Quantity: 2,
				Price:    10.0,
			},
		},
	}

	body, _ := json.Marshal(orderInput)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	apiInstance.CreateOrder(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var result dtos.OrderOutput
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, expectedOrder.ID, result.ID)
	mockCreateOrderUseCase.AssertExpectations(t)
}

func TestGetOrder(t *testing.T) {
	mockGetOrderUseCase := new(MockGetOrderUseCase)
	expectedOrder := dtos.OrderOutput{
		ID: "1",
	}

	mockGetOrderUseCase.On("Execute", mock.Anything, "1").Return(expectedOrder, nil)

	apiInstance := setupRouter(nil, nil, mockGetOrderUseCase, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	apiInstance.GetOrder(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var result dtos.OrderOutput
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, expectedOrder.ID, result.ID)
	mockGetOrderUseCase.AssertExpectations(t)
}

func TestListOrders(t *testing.T) {
	mockListOrderUseCase := new(MockListOrderUseCase)
	expectedList := dtos.ListOrderOutput{
		Orders: []dtos.OrderOutput{
			{ID: "1"},
			{ID: "2"},
		},
	}

	mockListOrderUseCase.On("Execute", mock.Anything).Return(expectedList, nil)

	apiInstance := setupRouter(nil, nil, nil, nil, mockListOrderUseCase)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	apiInstance.ListOrders(recorder, req)

	var result dtos.ListOrderOutput
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, len(expectedList.Orders), len(result.Orders))
	mockListOrderUseCase.AssertExpectations(t)
}

func TestUpdateOrder(t *testing.T) {
	mockUpdateOrderUseCase := new(MockUpdateOrderUseCase)
	expectedOrder := dtos.OrderOutput{
		ID: "1",
	}

	mockUpdateOrderUseCase.On("Execute", mock.Anything, "1", mock.Anything).Return(expectedOrder, nil)

	apiInstance := setupRouter(nil, mockUpdateOrderUseCase, nil, nil, nil)

	orderInput := dtos.OrderInput{
		CustomerName: "John Doe",
		Items: []dtos.ItemInput{
			{
				ID:       "1",
				Name:     "Item 1",
				Quantity: 2,
				Price:    10.0,
			},
		},
	}

	body, _ := json.Marshal(orderInput)

	req := httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	apiInstance.UpdateOrder(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var result dtos.OrderOutput
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, expectedOrder.ID, result.ID)
	mockUpdateOrderUseCase.AssertExpectations(t)
}

func TestCancelOrder(t *testing.T) {
	mockCancelOrderUseCase := new(MockCancelOrderUseCase)
	expectedOrder := dtos.OrderOutput{
		ID: "1",
	}

	mockCancelOrderUseCase.On("Execute", mock.Anything, "1").Return(expectedOrder, nil)

	apiInstance := setupRouter(nil, nil, mockCancelOrderUseCase, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	apiInstance.CancelOrder(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var result dtos.OrderOutput
	json.NewDecoder(recorder.Body).Decode(&result)

	assert.Equal(t, expectedOrder.ID, result.ID)
	mockCancelOrderUseCase.AssertExpectations(t)
}
