package models

import (
	"time"

	"github.com/google/uuid"
)

type Promotion struct {
	ID             uuid.UUID
	Title          string
	Description    string
	DiscountType   string
	DiscountValue  float64
	MinOrderAmount float64
	ValidFrom      time.Time
	ValidUntil     time.Time
	MaxUses        int32
	UsedCount      int32
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
type CreatePromotionRequest struct {
	Title          string    `json:"title" binding:"required"`
	Description    *string   `json:"description"`
	DiscountType   string    `json:"discount_type" binding:"required,oneof=percentage fixed buy_one_get_one"`
	DiscountValue  *float64  `json:"discount_value"`
	MinOrderAmount *float64  `json:"min_order_amount"`
	ValidFrom      time.Time `json:"valid_from" binding:"required"`
	ValidUntil     time.Time `json:"valid_until" binding:"required"`
	MaxUses        *int32    `json:"max_uses"`
	IsActive       *bool     `json:"is_active"`
}
type UpdatePromotionRequest struct {
    Title          *string    `json:"title"`
    Description    *string    `json:"description"`
    DiscountType   *string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed buy_one_get_one"`
    DiscountValue  *float64   `json:"discount_value"`
    MinOrderAmount *float64   `json:"min_order_amount"`
    ValidFrom      *time.Time `json:"valid_from"`
    ValidUntil     *time.Time `json:"valid_until"`
    MaxUses        *int32     `json:"max_uses"`
    IsActive       *bool      `json:"is_active"`
}

