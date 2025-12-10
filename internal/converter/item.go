package converter

import (
	"fmt"
	"math"
	"time"

	"github.com/kstsm/wb-sales-tracker/internal/dto"
	"github.com/kstsm/wb-sales-tracker/internal/models"
)

func ItemToResponse(item *models.Item) dto.ItemResponse {
	amountFloat := float64(item.Amount) / 100.0
	totalKopeks := int64(math.Round(amountFloat * 100))
	rubles := totalKopeks / 100
	kopeks := totalKopeks % 100
	amountStr := fmt.Sprintf("%d.%02d", rubles, kopeks)

	return dto.ItemResponse{
		ID:        item.ID.String(),
		Type:      item.Type,
		Amount:    amountStr,
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
