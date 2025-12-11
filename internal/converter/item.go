package converter

import (
	"fmt"
	"time"

	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
)

const (
	kopeksPerRuble = 100
)

func FormatRublesAmount(amount int) string {
	rubles := amount / kopeksPerRuble
	kopeks := amount % kopeksPerRuble

	if kopeks < 0 {
		kopeks = -kopeks
	}

	return fmt.Sprintf("%d.%02d", rubles, kopeks)
}

func ItemToResponse(item *models.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:        item.ID.String(),
		Type:      item.Type,
		Amount:    FormatRublesAmount(item.Amount),
		Date:      item.Date.UTC().Format(time.RFC3339),
		Category:  item.Category,
		CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func ItemsToResponse(items []*models.Item) []dto.ItemResponse {
	res := make([]dto.ItemResponse, len(items))
	for i, it := range items {
		res[i] = ItemToResponse(it)
	}

	return res
}
