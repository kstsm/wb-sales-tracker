package dto

import (
	"fmt"
	"strings"
	"time"
)

func buildWhereClause(from, to *time.Time, itemType, category *string) (string, []any) {
	var cond []string
	var args []any

	add := func(query string, val any) {
		cond = append(cond, fmt.Sprintf(query, len(args)+1))
		args = append(args, val)
	}

	if from != nil {
		add("date >= $%d", *from)
	}
	if to != nil {
		add("date <= $%d", *to)
	}
	if itemType != nil {
		add("type = $%d", *itemType)
	}
	if category != nil {
		add("category = $%d", *category)
	}

	if len(cond) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(cond, " AND "), args
}

func (r *GetItemsRequest) BuildOrder() string {
	sortBy := "date"
	sortOrder := "desc"

	if r.SortBy != nil {
		allowedSortBy := map[string]string{
			"date":     "date",
			"amount":   "amount",
			"category": "category",
		}
		if allowed, ok := allowedSortBy[strings.ToLower(*r.SortBy)]; ok {
			sortBy = allowed
		}
	}
	if r.SortOrder != nil {
		order := strings.ToUpper(*r.SortOrder)
		if order == "ASC" || order == "DESC" {
			sortOrder = order
		}
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}

func (r *AnalyticsRequest) GetGroupBy() string {
	if r.GroupBy == nil {
		return ""
	}
	allowed := map[string]string{
		"day":      "day",
		"week":     "week",
		"category": "category",
	}
	if val, ok := allowed[strings.ToLower(*r.GroupBy)]; ok {
		return val
	}
	return ""
}

func (r *AnalyticsRequest) BuildWhere() (string, []any) {
	var cond []string
	var args []any

	add := func(query string, val any) {
		cond = append(cond, fmt.Sprintf(query, len(args)+1))
		args = append(args, val)
	}

	if r.From != nil {
		fromStartOfDay := time.Date(r.From.Year(),
			r.From.Month(),
			r.From.Day(), 0, 0, 0, 0,
			r.From.Location())
		add("date >= $%d", fromStartOfDay)
	}
	if r.To != nil {
		toEndOfDay := time.Date(r.To.Year(),
			r.To.Month(),
			r.To.Day(), 23, 59, 59, 999999999,
			r.To.Location())
		add("date <= $%d", toEndOfDay)
	}

	if len(cond) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(cond, " AND "), args
}

func (r *GetItemsRequest) BuildWhere() (string, []any) {
	return buildWhereClause(r.From, r.To, r.Type, r.Category)
}
