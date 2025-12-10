package dto

import (
	"time"
)

type CreateItemRequest struct {
	Type     string `json:"type"      validate:"required,item_type"`
	Amount   int64  `json:"amount"    validate:"required"`
	Date     string `json:"date"      validate:"required,rfc3339"`
	Category string `json:"category"  validate:"required,min=3"`
}

type GetItemsRequest struct {
	From      *time.Time `json:"from,omitempty"`
	To        *time.Time `json:"to,omitempty"`
	Type      *string    `json:"type,omitempty" validate:"omitempty,item_type"`
	Category  *string    `json:"category,omitempty"`
	SortBy    *string    `json:"sort_by,omitempty" validate:"omitempty,sort_by"`
	SortOrder *string    `json:"sort_order,omitempty" validate:"omitempty,sort_order"`
}

type UpdateItemRequest struct {
	Type     *string    `json:"type,omitempty" validate:"omitempty,item_type"`
	Amount   *int64     `json:"amount,omitempty" validate:"omitempty"`
	Date     *time.Time `json:"date,omitempty" validate:"omitempty,rfc3339"`
	Category *string    `json:"category,omitempty" validate:"omitempty,min=3"`
}

type UpdateItemRequestInput struct {
	Type     *string `json:"type,omitempty" validate:"omitempty,item_type"`
	Amount   *int64  `json:"amount,omitempty" validate:"omitempty"`
	Date     *string `json:"date,omitempty" validate:"omitempty,rfc3339"`
	Category *string `json:"category,omitempty" validate:"omitempty,min=3"`
}

type AnalyticsRequest struct {
	From    *time.Time `json:"from,omitempty"`
	To      *time.Time `json:"to,omitempty"`
	GroupBy *string    `json:"group_by,omitempty"`
}
