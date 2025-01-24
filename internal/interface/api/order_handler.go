package api

import (
	"encoding/json"
	"net/http"

	"order-service/internal/application/dtos"

	"github.com/go-chi/chi"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func (api *API) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input dtos.OrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if input.CustomerName == "" || len(input.Items) == 0 {
		respondWithError(w, http.StatusBadRequest, "Customer name and items are required")
		return
	}

	orderOutput, err := api.createOrderUseCase.Execute(r.Context(), input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create order: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, orderOutput)
}

func (api *API) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderOutput, err := api.getOrderUseCase.Execute(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, orderOutput)
}

func (api *API) ListOrders(w http.ResponseWriter, r *http.Request) {
	listOutput, err := api.listOrderUseCase.Execute(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list orders")
		return
	}

	respondWithJSON(w, http.StatusOK, listOutput)
}

func (api *API) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input dtos.OrderInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	if input.CustomerName == "" || len(input.Items) == 0 {
		respondWithError(w, http.StatusBadRequest, "Customer name and items are required")
		return
	}

	orderOutput, err := api.updateOrderUseCase.Execute(r.Context(), id, input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	respondWithJSON(w, http.StatusOK, orderOutput)
}

func (api *API) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orderOutput, err := api.cancelOrderUseCase.Execute(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	respondWithJSON(w, http.StatusOK, orderOutput)
}
