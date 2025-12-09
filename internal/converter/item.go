package converter

import (
	"time"

	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
)

func ItemToResponse(item *models.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:        item.ID.String(),
		Type:      item.Type,
		Amount:    item.Amount,
		Date:      item.Date.Format(time.RFC3339),
		Category:  item.Category,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}
}

func ItemsToResponse(items []*models.Item) []dto.ItemResponse {
	res := make([]dto.ItemResponse, len(items))
	for i, it := range items {
		res[i] = ItemToResponse(it)
	}

	return res
}
