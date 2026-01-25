package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	OrderID    uuid.UUID
	MenuItemID uuid.UUID
	Rating     int32
	Comment    string
	Images     []string
	IsApproved bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type CreateReviewRequest struct {
	UserID     uuid.UUID   `json:"user_id" binding:"required"` // optional: from auth context
	OrderID    uuid.UUID   `json:"order_id" binding:"required"`
	MenuItemID uuid.UUID   `json:"menu_item_id" binding:"required"`
	Rating     *int32    `json:"rating" binding:"required,min=1,max=5"`
	Comment    *string  `json:"comment"`
	Images     []string `json:"images"`
}
type UpdateReviewRequest struct {
	Rating     *int32   `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment    *string  `json:"comment"`
	Images     []string `json:"images"`
}
