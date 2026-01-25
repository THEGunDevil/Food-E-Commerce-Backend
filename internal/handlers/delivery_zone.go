package handlers

import (
	"errors"
	"net/http"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateDeliveryZoneHandler(c *gin.Context) {
	var req models.CreateDeliveryZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	if req.AreaNames == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "area names is required",
		})
		return
	}
	if req.DeliveryFee == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "delivery fee is required",
		})
		return
	}
	if req.MinDeliveryTime == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "min delivery time is required",
		})
	}
	if req.MaxDeliveryTime == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "max delivery time is required",
		})
	}

	params := gen.CreateDeliveryZoneParams{
		ZoneName:        req.ZoneName,
		AreaNames:       req.AreaNames,
		DeliveryFee:     services.FloatToPGNumeric(req.DeliveryFee),
		MinDeliveryTime: pgtype.Int4{Int32: int32(req.MinDeliveryTime), Valid: true},
		MaxDeliveryTime: pgtype.Int4{Int32: int32(req.MaxDeliveryTime), Valid: true},
		IsActive:        services.PgTypeBool(req.IsActive),
	}
	res, err := db.Q.CreateDeliveryZone(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create delivery zone"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zone created successfully",
		Data:    services.ToDeliveryZoneResponse(res),
	})
}
func ListDeliveryZoneHandler(c *gin.Context) {
	res, err := db.Q.ListDeliveryZones(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list delivery zones"})
		return
	}
	var deliveryZones []models.DeliveryZone
	for _, zone := range res {
		deliveryZones = append(deliveryZones, services.ToDeliveryZoneResponse(zone))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zones listed successfully",
		Data:    deliveryZones,
	})
}
func ListActiveDeliveryZoneHandler(c *gin.Context) {
	res, err := db.Q.ListActiveDeliveryZones(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active delivery zones"})
		return
	}
	var deliveryZones []models.DeliveryZone
	for _, zone := range res {
		deliveryZones = append(deliveryZones, services.ToDeliveryZoneResponse(zone))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "active delivery zones listed successfully",
		Data:    deliveryZones,
	})
}
func GetDeliveryZoneByAreaHandler(c *gin.Context) {
	area := c.Param("area")
	if area == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "area is required"})
		return
	}
	res, err := db.Q.GetDeliveryZoneByArea(c, []string{area})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get delivery zone by area"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zone by area retrieved successfully",
		Data:    services.ToDeliveryZoneResponse(res),
	})
}
func UpdateDeliveryZoneHandler(c *gin.Context) {
	id := c.Param("delivery_zone_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delivery zone id is required"})
		return
	}

	var req models.UpdateDeliveryZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := gen.UpdateDeliveryZoneParams{
		ID: services.StringToPGUUID(id),
	}

	if req.ZoneName != nil {
		params.ZoneName = *req.ZoneName
	}

	if req.AreaNames != nil {
		params.AreaNames = req.AreaNames
	}

	if req.DeliveryFee != nil {
		params.DeliveryFee = services.FloatToPGNumeric(*req.DeliveryFee)
	}

	if req.MinDeliveryTime != nil {
		params.MinDeliveryTime = pgtype.Int4{Int32: *req.MinDeliveryTime, Valid: true}
	}

	if req.MaxDeliveryTime != nil {
		params.MaxDeliveryTime = pgtype.Int4{Int32: *req.MaxDeliveryTime, Valid: true}
	}

	if req.IsActive != nil {
		params.IsActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	res, err := db.Q.UpdateDeliveryZone(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "delivery zone not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to update delivery zone",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zone updated successfully",
		Data:    services.ToDeliveryZoneResponse(res),
	})
}

func ToggleDeliveryZoneStatusHandler(c *gin.Context) {
	id := c.Param("delivery_zone_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delivery zone id is required"})
		return
	}
	var req models.UpdateDeliveryZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var isActive bool
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	params := gen.ToggleDeliveryZoneStatusParams{
		ID:       services.StringToPGUUID(id),
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	}
	res, err := db.Q.ToggleDeliveryZoneStatus(c, params)
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "delivery zone not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to toggle delivery zone status"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zone status toggled successfully",
		Data:    services.ToDeliveryZoneResponse(res),
	})
}
func DeleteDeliveryZoneHandler(c *gin.Context) {
	id := c.Param("delivery_zone_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delivery zone id is required"})
		return
	}
	err := db.Q.DeleteDeliveryZone(c, services.StringToPGUUID(id))
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "delivery zone not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete delivery zone"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "delivery zone deleted successfully",
	})
}
