package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-sales-tracker/internal/middleware"
	"github.com/kstsm/wb-sales-tracker/internal/service"
)

type ItemManager interface {
	NewRouter() http.Handler
}

type Handler struct {
	service service.ItemManager
	log     *slog.Logger
	valid   *validator.Validate
}

func NewHandler(service service.ItemManager, log *slog.Logger, valid *validator.Validate) ItemManager {
	return &Handler{
		service: service,
		log:     log,
		valid:   valid,
	}
}

func (h *Handler) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CORS)

	r.Get("/", h.serveHTML("index.html"))

	h.registerAPIRoutes(r)

	return r
}

func (h *Handler) serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/"+filename)
	}
}
