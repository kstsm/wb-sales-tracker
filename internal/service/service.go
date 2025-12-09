package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
	"github.com/kstsm/wb-sales-tracker/internal/repository"
)

type ItemManager interface {
	CreateItem(ctx context.Context, req dto.CreateItemRequest) (*models.Item, error)
	GetItems(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, int, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequestInput) (*models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetAnalytics(ctx context.Context, req dto.AnalyticsRequest) (*dto.AnalyticsResponse, error)
	GetItemsForExport(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, error)
	ExportItemsCSV(ctx context.Context, req dto.GetItemsRequest) ([]byte, error)
}

type Service struct {
	repo repository.ItemManager
	log  *slog.Logger
}

func NewService(repo repository.ItemManager, log *slog.Logger) ItemManager {
	return &Service{
		repo: repo,
		log:  log,
	}
}
