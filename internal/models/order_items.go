package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID                  uuid.UUID
	OrderID             uuid.UUID
	MenuItemID          uuid.UUID
	MenuItemName        string
	Quantity            int32
	UnitPrice           float64
	TotalPrice          float64
	SpecialInstructions *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
type CreateOrderItemRequest struct {
	MenuItemID          *string `json:"menu_item_id"` // optional: menu deleted?
	MenuItemName        string  `json:"menu_item_name" binding:"required"`
	Quantity            int32   `json:"quantity" binding:"required,gt=0"`
	UnitPrice           float64 `json:"unit_price" binding:"required,gt=0"`
	TotalPrice          float64 `json:"total_price" binding:"required,gt=0"`
	SpecialInstructions *string `json:"special_instructions"`
}
type UpdateOrderItemRequest struct {
	MenuItemID          *string  `json:"menu_item_id"`
	MenuItemName        *string  `json:"menu_item_name"`
	Quantity            *int32   `json:"quantity" binding:"omitempty,gt=0"`
	UnitPrice           *float64 `json:"unit_price" binding:"omitempty,gt=0"`
	TotalPrice          *float64 `json:"total_price" binding:"omitempty,gt=0"`
	SpecialInstructions *string  `json:"special_instructions"`
}
