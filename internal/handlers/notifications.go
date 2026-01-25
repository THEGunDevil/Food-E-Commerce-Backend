package handlers

import (
	"errors"
	"net/http"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetNotificationByIDHandler(c *gin.Context) {
	idStr := c.Param("notification_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "notification id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid notification id",
		})
		return
	}
	notification, err := db.Q.GetNotificationByID(c, services.UUIDToPGType(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to get notification",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "notification retrieved successfully",
		Data:    services.ToNotificationResponse(notification),
	})
}

func ListNotificationsHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "user id is required",
		})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid user id",
		})
		return
	}
	notifications, err := db.Q.ListNotificationsByUser(c, services.UUIDToPGType(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "notifications not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list notifications",
		})
		return
	}

	res := make([]models.Notification, 0, len(notifications))
	for _, notification := range notifications {
		res = append(res, services.ToNotificationResponse(notification))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "notifications retrieved successfully",
		Data:    res,
	})
}

func MarkNotificationAsReadHandler(c *gin.Context) {
	idStr := c.Param("notification_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "notification id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid notification id",
		})
		return
	}
	n,err := db.Q.MarkNotificationRead(c, services.UUIDToPGType(id))
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to mark notification as read",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "notification marked as read successfully",
		Data:    services.ToNotificationResponse(n),
	}) 
}

func DeleteNotificationByIDHandler(c *gin.Context) {
	idStr := c.Param("notification_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "notification id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid notification id",
		})
		return
	}
	err = db.Q.DeleteNotification(c, services.UUIDToPGType(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete notification",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "notification deleted successfully",
	})
}
