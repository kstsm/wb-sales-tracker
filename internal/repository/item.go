package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-sales-tracker/internal/apperrors"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
	"github.com/kstsm/wb-sales-tracker/internal/repository/queries"
)

func (r *Repository) CreateItem(ctx context.Context, item models.Item) error {
	amountFloat := float64(item.Amount) / 100.0

	_, err := r.conn.Exec(ctx, queries.CreateItemQuery,
		item.ID,
		item.Type,
		amountFloat,
		item.Date,
		item.Category,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("Exec-CreateItem: %w", err)
	}

	return nil
}

func (r *Repository) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	var item models.Item
	var amountFloat float64

	err := r.conn.QueryRow(ctx, queries.GetItemByIDQuery, id).Scan(
		&item.ID,
		&item.Type,
		&amountFloat,
		&item.Date,
		&item.Category,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("QueryRow-getItemByID: %w", err)
	}

	item.Amount = int64(amountFloat * 100)

	return &item, nil
}

func (r *Repository) GetItems(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, int, error) {
	whereClause, args := r.buildItemsWhere(req)
	orderClause := r.buildItemsOrder(req)

	var total int
	if err := r.conn.QueryRow(ctx, queries.BaseCountQuery+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("QueryRow-GetItems: %w", err)
	}

	rows, err := r.conn.Query(ctx, queries.BaseSelectQuery+whereClause+orderClause, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Query-GetItems: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		var amountFloat float64
		if err = rows.Scan(
			&item.ID,
			&item.Type,
			&amountFloat,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("Scan-GetItems: %w", err)
		}
		item.Amount = int64(amountFloat * 100)
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Err-GetItems: %w", err)
	}

	return items, total, nil
}

func (r *Repository) UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequest) (*models.Item, error) {
	var item models.Item
	var amountFloat float64

	var amountVal interface{}
	if req.Amount != nil {
		amountVal = float64(*req.Amount) / 100.0
	} else {
		amountVal = nil
	}

	err := r.conn.QueryRow(
		ctx,
		queries.UpdateItemQuery,
		id,
		req.Type,
		amountVal,
		req.Date,
		req.Category,
	).Scan(
		&item.ID,
		&item.Type,
		&amountFloat,
		&item.Date,
		&item.Category,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrItemNotFound
		}
		return nil, fmt.Errorf("QueryRow-UpdateItem: %w", err)
	}

	item.Amount = int64(amountFloat * 100)

	return &item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	var deletedID uuid.UUID
	err := r.conn.QueryRow(ctx, queries.DeleteItemQuery, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrItemNotFound
		}
		return fmt.Errorf("QueryRow-DeleteItem: %w", err)
	}

	return nil
}

func (r *Repository) GetItemsForExport(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, error) {
	whereClause, args := r.buildItemsWhere(req)
	orderClause := r.buildItemsOrder(req)

	rows, err := r.conn.Query(ctx, queries.BaseSelectQuery+whereClause+orderClause, args...)
	if err != nil {
		return nil, fmt.Errorf("Query-GetItemsForExport: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		var amountFloat float64
		if err = rows.Scan(
			&item.ID,
			&item.Type,
			&amountFloat,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan-GetItemsForExport: %w", err)
		}
		item.Amount = int64(amountFloat * 100)
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Err-GetItemsForExport: %w", err)
	}

	return items, nil
}

func (r *Repository) buildItemsWhere(req dto.GetItemsRequest) (string, []any) {
	var cond []string
	var args []any

	add := func(query string, val any) {
		cond = append(cond, fmt.Sprintf(query, len(args)+1))
		args = append(args, val)
	}

	if req.From != nil {
		add("date >= $%d", *req.From)
	}
	if req.To != nil {
		add("date <= $%d", *req.To)
	}
	if req.Type != nil {
		add("type = $%d", *req.Type)
	}
	if req.Category != nil {
		add("category = $%d", *req.Category)
	}

	if len(cond) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(cond, " AND "), args
}

func (r *Repository) buildItemsOrder(req dto.GetItemsRequest) string {
	sortBy := "date"
	sortOrder := "desc"

	if req.SortBy != nil {
		allowedSortBy := map[string]string{
			"date":     "date",
			"amount":   "amount",
			"category": "category",
		}
		if allowed, ok := allowedSortBy[strings.ToLower(*req.SortBy)]; ok {
			sortBy = allowed
		}
	}
	if req.SortOrder != nil {
		order := strings.ToUpper(*req.SortOrder)
		if order == "ASC" || order == "DESC" {
			sortOrder = order
		}
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
}
