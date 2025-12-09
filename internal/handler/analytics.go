package handler

import (
	"net/http"

	"github.com/kstsm/wb-sales-tracker/internal/dto"
)

func (h *Handler) getAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AnalyticsRequest

	if err := parseAnalyticsQuery(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetAnalytics(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}
