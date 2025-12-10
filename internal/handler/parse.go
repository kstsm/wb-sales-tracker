package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kstsm/wb-sales-tracker/internal/apperrors"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
)

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

	allowedParams := map[string]bool{
		"from":       true,
		"to":         true,
		"type":       true,
		"category":   true,
		"sort_by":    true,
		"sort_order": true,
	}

	for param := range q {
		if !allowedParams[param] {
			return fmt.Errorf("unknown parameter '%s'", param)
		}
	}

	var err error
	if fromStr := q.Get("from"); fromStr != "" {
		req.From, err = parseDate(fromStr)
		if err != nil && !errors.Is(err, apperrors.ErrEmptyDate) {
			return err
		}
	}

	if toStr := q.Get("to"); toStr != "" {
		req.To, err = parseDate(toStr)
		if err != nil && !errors.Is(err, apperrors.ErrEmptyDate) {
			return err
		}
	}

	if req.From != nil && req.To != nil && req.From.After(*req.To) {
		return errors.New("parameter 'from' cannot be after 'to'")
	}

	typeStr := strings.TrimSpace(q.Get("type"))
	if typeStr != "" {
		req.Type = &typeStr
	} else if q.Has("type") {
		return errors.New("parameter 'type' cannot be empty")
	}

	categoryStr := strings.TrimSpace(q.Get("category"))
	if categoryStr != "" {
		req.Category = &categoryStr
	}

	sortByStr := strings.TrimSpace(q.Get("sort_by"))
	if sortByStr != "" {
		req.SortBy = &sortByStr
	} else if q.Has("sort_by") {
		return errors.New("parameter 'sort_by' cannot be empty")
	}

	sortOrderStr := strings.TrimSpace(q.Get("sort_order"))
	if sortOrderStr != "" {
		req.SortOrder = &sortOrderStr
	} else if q.Has("sort_order") {
		return errors.New("parameter 'sort_order' cannot be empty")
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

func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, apperrors.ErrEmptyDate
	}

	s = strings.TrimSpace(s)
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
