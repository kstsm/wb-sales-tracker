package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
)

type ItemManager interface {
	CreateItem(ctx context.Context, item models.Item) error
	GetItemByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	GetItems(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, int, error)
	UpdateItem(ctx context.Context, id uuid.UUID, req dto.UpdateItemRequest) (*models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetItemsForExport(ctx context.Context, req dto.GetItemsRequest) ([]*models.Item, error)
	GetAnalytics(ctx context.Context, req dto.AnalyticsRequest) (*dto.AnalyticsResponse, error)
}

type Repository struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func NewRepository(conn *pgxpool.Pool, log *slog.Logger) ItemManager {
	return &Repository{
		conn: conn,
		log:  log,
	}
}
