package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kstsm/wb-sales-tracker/internal/apperrors"
	"github.com/kstsm/wb-sales-tracker/internal/converter"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
)

func (h *Handler) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.CreateItem(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusCreated, resp)
}

func (h *Handler) getItemByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetItemByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "Item not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "Failed to get item")
		}
		return
	}

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusOK, resp)
}

func (h *Handler) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetItemsRequest

	if err := parseGetItemsQuery(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, total, err := h.service.GetItems(r.Context(), req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := converter.ItemsToResponse(result)
	h.respondJSON(w, http.StatusOK, dto.ItemsListResponse{
		Items: resp,
		Total: total,
	})
}

func (h *Handler) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateItemRequestInput
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err = h.valid.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.UpdateItem(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "Item not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := converter.ItemToResponse(result)
	h.respondJSON(w, http.StatusOK, resp)
}

func (h *Handler) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteItem(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "Item not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.respondJSON(w, http.StatusOK, nil)
}

func (h *Handler) exportItemCSVHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetItemsRequest

	if err := parseGetItemsQuery(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.ExportItemsCSV(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrItemNotFound):
			h.respondError(w, http.StatusNotFound, "Item not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.respondCSV(w, http.StatusOK, data)
}
