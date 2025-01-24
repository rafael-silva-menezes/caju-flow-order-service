package api

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(api *API) *chi.Mux {
	r := chi.NewRouter()
	RegisterMiddlewares(r)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", api.CreateOrder)
		r.Get("/", api.ListOrders)
		r.Get("/{id}", api.GetOrder)
		r.Put("/{id}", api.UpdateOrder)
		r.Delete("/{id}", api.CancelOrder)
	})

	return r
}
