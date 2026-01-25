package models

import (
	"time"

	"github.com/google/uuid"
)

type FavoriteResponse struct {
	ID         uuid.UUID `json:"id" binding:"required,uuid4"`
	UserID     uuid.UUID `json:"user_id" binding:"required,uuid4"`
	MenuItemID uuid.UUID `json:"menu_id" binding:"required,uuid4"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateFavoriteRequest struct {
	UserID     uuid.UUID `json:"user_id" binding:"required,uuid4"`
	MenuItemID uuid.UUID `json:"menu_id" binding:"required,uuid4"`
}
