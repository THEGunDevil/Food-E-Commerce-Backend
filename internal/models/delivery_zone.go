package models

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryZone struct {
	ID              uuid.UUID
	ZoneName        string
	AreaNames       []string
	DeliveryFee     float64
	MinDeliveryTime int
	MaxDeliveryTime int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
type CreateDeliveryZoneRequest struct {
    ZoneName        string   `json:"zone_name" binding:"required"`
    AreaNames       []string `json:"area_names" binding:"required,min=1"`
    DeliveryFee     *float64  `json:"delivery_fee" binding:"required"`
    MinDeliveryTime int   `json:"min_delivery_time"`
    MaxDeliveryTime int   `json:"max_delivery_time"`
    IsActive        *bool    `json:"is_active"`
}
type UpdateDeliveryZoneRequest struct {
    ZoneName        *string  `json:"zone_name"`
    AreaNames       []string `json:"area_names"`
    DeliveryFee     *float64 `json:"delivery_fee"`
    MinDeliveryTime *int32   `json:"min_delivery_time"`
    MaxDeliveryTime *int32   `json:"max_delivery_time"`
    IsActive        *bool    `json:"is_active"`
}


