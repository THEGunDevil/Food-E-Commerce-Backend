package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
)

func Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("notification dispatcher stopped")
			return
		case <-ticker.C:
			process(ctx)
		}
	}
}
func process(ctx context.Context) {
	events, err := db.Q.ListUndeliveredEvents(ctx)
	if err != nil {
		log.Println("failed to fetch events:", err)
		return
	}

	for _, event := range events {
		if err := handleEvent(ctx, models.Event{
			ID:          event.ID.Bytes,
			EventType:   event.EventType,
			Payload:     event.Payload,
			Delivered:   event.Delivered.Bool,
			CreatedAt:   event.CreatedAt.Time,
			DeliveredAt: &event.DeliveredAt.Time,
		}); err != nil {
			log.Println("event handling failed:", err)
			continue
		}

		_, err := db.Q.MarkEventDelivered(ctx, event.ID)
		if err != nil {
			log.Println("failed to mark event delivered:", err)
		}
	}
}

func handleEvent(ctx context.Context, event models.Event) error {
	switch event.EventType {

	case "order":
		return handleOrderEvent(ctx, event)
	default:
		// Ignore unknown events safely
		return nil
	}
}

func handleOrderEvent(ctx context.Context, event models.Event) error {
	var payload struct {
		UserID  string `json:"user_id"`
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	_, err := db.Q.CreateNotification(ctx, gen.CreateNotificationParams{
		UserID: StringToPGUUID(payload.UserID),
		Title:  "Order Update",
		Message:   "Your order " + payload.OrderID + " is now " + payload.Status,
	})

	return err
}