package handlers

import (
	"errors"
	"net/http"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateEventHandler(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	event, err := db.Q.CreateEvent(c, gen.CreateEventParams{
		EventType: req.EventType,
		Payload:   req.Payload,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "event created successfully",
		Data:    services.ToEventResponse(event),
	})
}
func ListUndeliveredEventsHandler(c *gin.Context) {
	events, err := db.Q.ListUndeliveredEvents(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, models.APIResponse{
				Success: true,
				Message: "no undelivered events",
				Data:    []models.Event{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list events",
		})
		return
	}

	res := make([]models.Event, 0, len(events))
	for _, e := range events {
		res = append(res, services.ToEventResponse(e))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "undelivered events retrieved successfully",
		Data:    res,
	})
}
func MarkEventDeliveredHandler(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "event id is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid event id",
		})
		return
	}

	event, err := db.Q.MarkEventDelivered(c, services.UUIDToPGType(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "event not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to mark event as delivered",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "event marked as delivered",
		Data:    services.ToEventResponse(event),
	})
}
func DeleteEventHandler(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "event id is required",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid event id",
		})
		return
	}

	if err := db.Q.DeleteEvent(c, services.UUIDToPGType(id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "event not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete event",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "event deleted successfully",
	})
}
