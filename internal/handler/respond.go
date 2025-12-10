package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kstsm/wb-sales-tracker/internal/models"
)

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.log.Errorf("respondJSON: %v", err.Error())
		return
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(models.Error{Error: message})
	if err != nil {
		h.log.Errorf("respondError: %v", err.Error())
		return
	}
}

func (h *Handler) respondCSV(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=items.csv")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		h.log.Errorf("respondJSON: %v", err.Error())
		return
	}
}
