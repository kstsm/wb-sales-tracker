package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/wb-sales-tracker/internal/apperrors"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
	"github.com/kstsm/wb-sales-tracker/internal/repository/queries"
)

func (r *Repository) CreateItem(ctx context.Context, item models.Item) error {
	_, err := r.conn.Exec(ctx, queries.CreateItemQuery,
		item.ID,
		item.Type,
		item.Amount,
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

	err := r.conn.QueryRow(ctx, queries.GetItemByIDQuery, id).Scan(
		&item.ID,
		&item.Type,
		&item.Amount,
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

	return &item, nil
}

func (r *Repository) GetItems(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, int, error) {
	whereClause, args := req.BuildWhere()
	orderClause := req.BuildOrder()

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
		if err = rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("Scan-GetItems: %w", err)
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("Err-GetItems: %w", err)
	}

	return items, total, nil
}

func (r *Repository) UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequest) (*models.Item, error) {
	var item models.Item

	var amountVal interface{}
	if req.Amount != nil {
		amountVal = req.Amount
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
		&item.Amount,
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

	return &item, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.Exec(ctx, queries.DeleteItemQuery, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrItemNotFound
		}
		return fmt.Errorf("QueryRow-DeleteItem: %w", err)
	}

	return nil
}

func (r *Repository) GetItemsForExport(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, error) {
	whereClause, args := req.BuildWhere()
	orderClause := req.BuildOrder()

	rows, err := r.conn.Query(ctx, queries.BaseSelectQuery+whereClause+orderClause, args...)
	if err != nil {
		return nil, fmt.Errorf("Query-GetItemsForExport: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err = rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Date,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan-GetItemsForExport: %w", err)
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Err-GetItemsForExport: %w", err)
	}

	return items, nil
}
