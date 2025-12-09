package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kstsm/wb-sales-tracker/internal/converter"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
	"github.com/kstsm/wb-sales-tracker/pkg/export"
)

func (s *Service) CreateItem(ctx context.Context, req dto.CreateItemRequest) (*models.Item, error) {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	item := models.Item{
		ID:        uuid.New(),
		Type:      req.Type,
		Amount:    req.Amount,
		Date:      date,
		Category:  req.Category,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err = s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *Service) ExportItemsCSV(ctx context.Context, req dto.GetItemsRequest) ([]byte, error) {
	items, err := s.repo.GetItemsForExport(ctx, req)
	if err != nil {
		return nil, err
	}

	exportItems := converter.ItemsToResponse(items)

	var buf bytes.Buffer
	if writeErr := export.WriteItemsCSV(&buf, nil, exportItems); writeErr != nil {
		return nil, writeErr
	}

	return buf.Bytes(), nil
}

func (s *Service) UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequestInput) (*models.Item, error) {
	item := dto.UpdateItemRequest{
		Type:     req.Type,
		Amount:   req.Amount,
		Category: req.Category,
	}

	if req.Date != nil {
		date, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		item.Date = &date
	}

	return s.repo.UpdateItem(ctx, id, item)
}

func (s *Service) GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	return s.repo.GetItemByID(ctx, id)
}

func (s *Service) GetItems(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, int, error) {
	return s.repo.GetItems(ctx, req)
}

func (s *Service) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteItem(ctx, id)
}

func (s *Service) GetItemsForExport(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, error) {
	return s.repo.GetItemsForExport(ctx, req)
}
