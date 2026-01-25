package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAddress struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Label        string    `json:"label"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 *string   `json:"address_line2,omitempty"`
	Area         string    `json:"area"`
	City         string    `json:"city"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserAddressRequest struct {
	UserID       string   `json:"user_id" binding:"required,uuid"`
	Label        string   `json:"label" binding:"required"`
	AddressLine1 string   `json:"address_line1" binding:"required"`
	AddressLine2 *string  `json:"address_line2,omitempty"`
	Area         string   `json:"area" binding:"required"`
	City         *string  `json:"city,omitempty"` // optional, default 'Dhaka'
	PostalCode   *string  `json:"postal_code,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	IsDefault    *bool    `json:"is_default,omitempty"` // optional, default false
}

type UpdateUserAddressRequest struct {
	Label        *string  `json:"label,omitempty"`
	AddressLine1 *string  `json:"address_line1,omitempty"`
	AddressLine2 *string  `json:"address_line2,omitempty"`
	Area         *string  `json:"area,omitempty"`
	City         *string  `json:"city,omitempty"`
	PostalCode   *string  `json:"postal_code,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	IsDefault    *bool    `json:"is_default,omitempty"`
}