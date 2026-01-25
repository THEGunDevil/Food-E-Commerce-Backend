package models

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	EventID   uuid.UUID `json:"event_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Priority  string    `json:"priority"`
	IsRead    bool      `json:"is_read"`
	Metadata  []byte    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateNotificationRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	EventID  string `json:"event_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Message  string `json:"message" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=order promotion system review payment delivery"`
	Priority string `json:"priority" binding:"omitempty,oneof=high medium low"`
	Metadata []byte `json:"metadata"`
}
type UpdateNotificationRequest struct {
	IsRead bool `json:"is_read"`
}

