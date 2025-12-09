package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
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

func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	value := chi.URLParam(r, param)
	if strings.TrimSpace(value) == "" {
		return uuid.Nil, fmt.Errorf("%s is required", param)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s", param)
	}

	return id, nil
}

func parseGetItemsQuery(r *http.Request, req *dto.GetItemsRequest) error {
	q := r.URL.Query()

	var err error
	if fromStr := q.Get("from"); fromStr != "" {
		req.From, err = parseDate(fromStr)
		if err != nil && !errors.Is(err, errEmptyDate) {
			return err
		}
	}

	if toStr := q.Get("to"); toStr != "" {
		req.To, err = parseDate(toStr)
		if err != nil && !errors.Is(err, errEmptyDate) {
			return err
		}
	}

	if req.From != nil && req.To != nil && req.From.After(*req.To) {
		return errors.New("parameter 'from' cannot be after 'to'")
	}

	if typeStr := q.Get("type"); typeStr != "" {
		req.Type = &typeStr
	}

	if categoryStr := q.Get("category"); categoryStr != "" {
		req.Category = &categoryStr
	}

	if sortByStr := q.Get("sort_by"); sortByStr != "" {
		req.SortBy = &sortByStr
	}

	if sortOrderStr := q.Get("sort_order"); sortOrderStr != "" {
		req.SortOrder = &sortOrderStr
	}

	return nil
}

func parseAnalyticsQuery(r *http.Request, req *dto.AnalyticsRequest) error {
	q := r.URL.Query()

	fromStr := q.Get("from")
	if fromStr == "" {
		return errors.New("parameter 'from' is required")
	}

	toStr := q.Get("to")
	if toStr == "" {
		return errors.New("parameter 'to' is required")
	}

	var err error
	req.From, err = parseDate(fromStr)
	if err != nil {
		return fmt.Errorf("invalid 'from' date format, expected YYYY-MM-DD: %w", err)
	}

	req.To, err = parseDate(toStr)
	if err != nil {
		return fmt.Errorf("invalid 'to' date format, expected YYYY-MM-DD: %w", err)
	}

	if req.From.After(*req.To) {
		return errors.New("parameter 'from' cannot be after 'to'")
	}

	if groupByStr := q.Get("group_by"); groupByStr != "" {
		req.GroupBy = &groupByStr
	}

	return nil
}

var errEmptyDate = errors.New("empty date string")

func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, errEmptyDate
	}

	s = strings.TrimSpace(s)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
