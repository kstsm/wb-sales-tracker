package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) registerAPIRoutes(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/items", h.createItemHandler)
		r.Get("/items", h.getItemsHandler)
		r.Get("/items/{id}", h.getItemByIDHandler)
		r.Put("/items/{id}", h.updateItemHandler)
		r.Delete("/items/{id}", h.deleteItemHandler)
		r.Get("/analytics", h.getAnalyticsHandler)
		r.Get("/export", h.exportItemCSVHandler)
	})
}
