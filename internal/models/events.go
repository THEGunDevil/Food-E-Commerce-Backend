package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)
type Event struct {
	ID          uuid.UUID  `json:"id"`
	EventType   string     `json:"event_type"`
	Payload     json.RawMessage `json:"payload"`
	Delivered   bool       `json:"delivered"`
	CreatedAt   time.Time  `json:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

type CreateEventRequest struct {
	EventType string          `json:"event_type" binding:"required"`
	Payload   json.RawMessage `json:"payload" binding:"required"`
}