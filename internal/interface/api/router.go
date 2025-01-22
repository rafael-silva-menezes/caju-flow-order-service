package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func addMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

func addRoutes(r *chi.Mux, api *API) {
	r.Post("/orders", api.CreateOrder)
	r.Get("/orders/{id}", api.GetOrder)
	r.Get("/orders", api.ListOrders)
	r.Put("/orders/{id}", api.UpdateOrder)
	r.Delete("/orders/{id}", api.CancelOrder)
}

func NewRouter(api *API) *chi.Mux {
	r := chi.NewRouter()
	addMiddlewares(r)
	addRoutes(r, api)
	return r
}
