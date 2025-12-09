package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
