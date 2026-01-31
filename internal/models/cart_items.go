package models

import (
	"time"

	"github.com/google/uuid"
)

type CartItems struct {
	ID                  uuid.UUID `json:"id"`
	UserID              uuid.UUID `json:"user_id"`
	MenuItemID          uuid.UUID `json:"menu_item_id"`
	SessionID           uuid.UUID `json:"session_id"`
	Quantity            int       `json:"quantity"`
	SpecialInstructions string    `json:"special_instructions"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
type CartItemResponse struct {
	CartID     uuid.UUID `json:"cart_id"`
	MenuItemID uuid.UUID `json:"menu_item_id"`

	Name            string          `json:"name"`
	CategoryName    string          `json:"category_name"`
	Price           float64         `json:"price"`
	DiscountedPrice float64         `json:"discounted_price"`
	OriginalPrice   float64         `json:"original_price"`
	LineSubtotal    float64         `json:"line_subtotal"`
	Image           []MenuItemImage `json:"image"`

	Quantity            int    `json:"quantity"`
	SpecialInstructions string `json:"special_instructions"`
}

type CartItemsReq struct {
	MenuItemID          uuid.UUID `json:"menu_item_id" validate:"required"`
	Quantity            int       `json:"quantity" validate:"required, gte=1,lte=10"`
	SpecialInstructions *string   `json:"special_instructions"`
}
type UpdateCartItemQuantityReq struct {
	Quantity *int `json:"quantity" validate:"required, gte=1,lte=10"`
}
