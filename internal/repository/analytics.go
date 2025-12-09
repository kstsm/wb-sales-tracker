package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kstsm/wb-sales-tracker/internal/converter"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/repository/queries"
)

func (r *Repository) GetAnalytics(ctx context.Context, req dto.AnalyticsRequest) (*dto.AnalyticsResponse, error) {
	whereClause, args := req.BuildWhere()
	if groupBy := req.GetGroupBy(); groupBy != "" {
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

	sumValue := 0.0
	if sum.Valid {
		sumValue = sum.Float64
	}

	fromStr := req.From.Format("2006-01-02")
	toStr := req.To.Format("2006-01-02")

	return &dto.AnalyticsResponse{
		From:         fromStr,
		To:           toStr,
		Sum:          sumValue,
		Avg:          converter.ToFloatPtrIf(avg, count > 0),
		Count:        count,
		Median:       converter.ToFloatPtrIf(median, count > 0),
		Percentile90: converter.ToFloatPtrIf(percentile90, count > 0),
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

	for rows.Next() {
		var groupKey any
		var sum, avg, median, percentile90 sql.NullFloat64
		var count int

		if scanErr := rows.Scan(&groupKey, &sum, &avg, &count, &median, &percentile90); scanErr != nil {
			return nil, fmt.Errorf("Scan-GetGroupedAnalytics: %w", scanErr)
		}

		if sum.Valid {
			totalSum += sum.Float64
		}
		totalCount += count

		grouped = append(grouped, dto.GroupedAnalytics{
			Group:        converter.FormatGroupKey(groupKey, groupBy),
			Sum:          converter.ToFloatPtr(sum),
			Avg:          converter.ToFloatPtrIf(avg, count > 0),
			Count:        count,
			Median:       converter.ToFloatPtrIf(median, count > 0),
			Percentile90: converter.ToFloatPtrIf(percentile90, count > 0),
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

	fromStr := req.From.Format("2006-01-02")
	toStr := req.To.Format("2006-01-02")

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
