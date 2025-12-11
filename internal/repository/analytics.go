package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kstsm/wb-sales-tracker/internal/converter"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/repository/queries"
)

func (r *Repository) GetAnalytics(ctx context.Context, req dto.AnalyticsRequest) (*dto.AnalyticsResponse, error) {
	whereClause, args := r.buildAnalyticsWhere(req)
	if groupBy := r.getAnalyticsGroupBy(req); groupBy != "" {
		return r.getGroupedAnalytics(ctx, whereClause, args, groupBy, req)
	}

	var sum sql.NullFloat64
	var avg, median, percentile90 sql.NullFloat64
	var count int

	err := r.conn.QueryRow(ctx, fmt.Sprintf(queries.AnalyticsQuery, whereClause), args...).
		Scan(&sum, &avg, &count, &median, &percentile90)
	if err != nil {
		return nil, fmt.Errorf("QueryRow-GetAnalytics: %w", err)
	}

	const kopeksPerRuble = 100.0

	sumValue := 0.0
	if sum.Valid {
		sumValue = sum.Float64 / kopeksPerRuble
	}

	var avgValue *float64
	if avg.Valid && count > 0 {
		val := avg.Float64 / kopeksPerRuble
		avgValue = &val
	}

	var medianValue *float64
	if median.Valid && count > 0 {
		val := median.Float64 / kopeksPerRuble
		medianValue = &val
	}

	var percentile90Value *float64
	if percentile90.Valid && count > 0 {
		val := percentile90.Float64 / kopeksPerRuble
		percentile90Value = &val
	}

	fromStr := req.From.Format(time.RFC3339)
	toStr := req.To.Format(time.RFC3339)

	return &dto.AnalyticsResponse{
		From:         fromStr,
		To:           toStr,
		Sum:          sumValue,
		Avg:          avgValue,
		Count:        count,
		Median:       medianValue,
		Percentile90: percentile90Value,
	}, nil
}

func (r *Repository) getGroupedAnalytics(ctx context.Context,
	whereClause string,
	args []any,
	groupBy string,
	req dto.AnalyticsRequest,
) (*dto.AnalyticsResponse, error) {
	query, err := r.getGroupedQuery(groupBy, whereClause)
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Query-GetGroupedAnalytics: %w", err)
	}
	defer rows.Close()

	var grouped []dto.GroupedAnalytics
	var totalSum float64
	var totalCount int

	const kopeksPerRuble = 100.0

	for rows.Next() {
		var groupKey any
		var sum, avg, median, percentile90 sql.NullFloat64
		var count int

		if scanErr := rows.Scan(&groupKey, &sum, &avg, &count, &median, &percentile90); scanErr != nil {
			return nil, fmt.Errorf("Scan-GetGroupedAnalytics: %w", scanErr)
		}

		var sumValue *float64
		if sum.Valid {
			val := sum.Float64 / kopeksPerRuble
			sumValue = &val
			totalSum += val
		}

		var avgValue *float64
		if avg.Valid && count > 0 {
			val := avg.Float64 / kopeksPerRuble
			avgValue = &val
		}

		var medianValue *float64
		if median.Valid && count > 0 {
			val := median.Float64 / kopeksPerRuble
			medianValue = &val
		}

		var percentile90Value *float64
		if percentile90.Valid && count > 0 {
			val := percentile90.Float64 / kopeksPerRuble
			percentile90Value = &val
		}

		totalCount += count

		grouped = append(grouped, dto.GroupedAnalytics{
			Group:        converter.FormatGroupKey(groupKey, groupBy),
			Sum:          sumValue,
			Avg:          avgValue,
			Count:        count,
			Median:       medianValue,
			Percentile90: percentile90Value,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err-GetGroupedAnalytics: %w", err)
	}

	var totalAvg *float64
	if totalCount > 0 {
		avg := totalSum / float64(totalCount)
		totalAvg = &avg
	}

	fromStr := req.From.Format(time.RFC3339)
	toStr := req.To.Format(time.RFC3339)

	return &dto.AnalyticsResponse{
		From:    fromStr,
		To:      toStr,
		Sum:     totalSum,
		Avg:     totalAvg,
		Count:   totalCount,
		Grouped: grouped,
	}, nil
}

func (r *Repository) getGroupedQuery(groupBy, whereClause string) (string, error) {
	switch groupBy {
	case "day":
		return fmt.Sprintf(queries.AnalyticsGroupedByDayQuery, whereClause), nil
	case "week":
		return fmt.Sprintf(queries.AnalyticsGroupedByWeekQuery, whereClause), nil
	case "category":
		return fmt.Sprintf(queries.AnalyticsGroupedByCategoryQuery, whereClause), nil
	default:
		return "", fmt.Errorf("unsupported group_by value: %s", groupBy)
	}
}

func (r *Repository) buildAnalyticsWhere(req dto.AnalyticsRequest) (string, []any) {
	var cond []string
	var args []any

	add := func(query string, val any) {
		cond = append(cond, fmt.Sprintf(query, len(args)+1))
		args = append(args, val)
	}

	if req.From != nil {
		fromStartOfDay := time.Date(req.From.Year(),
			req.From.Month(),
			req.From.Day(), 0, 0, 0, 0,
			req.From.Location())
		add("date >= $%d", fromStartOfDay)
	}
	if req.To != nil {
		toEndOfDay := time.Date(req.To.Year(),
			req.To.Month(),
			req.To.Day(), 23, 59, 59, 999999999,
			req.To.Location())
		add("date <= $%d", toEndOfDay)
	}

	if len(cond) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(cond, " AND "), args
}

func (r *Repository) getAnalyticsGroupBy(req dto.AnalyticsRequest) string {
	if req.GroupBy == nil {
		return ""
	}
	allowed := map[string]string{
		"day":      "day",
		"week":     "week",
		"category": "category",
	}
	if val, ok := allowed[strings.ToLower(*req.GroupBy)]; ok {
		return val
	}
	return ""
}
