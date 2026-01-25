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

func CreateOrderItemHandler(c *gin.Context) {
	idString := c.Param("id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	menuItemID, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}
	var req models.CreateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "quantity must be greater than 0",
		})
		return
	}
	if req.UnitPrice <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "unit price must be greater than 0",
		})
		return
	}
	if req.TotalPrice <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "total price must be greater than 0",
		})
		return
	}
	if req.MenuItemName == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu item name is required",
		})
		return
	}
	if req.SpecialInstructions == nil {
		req.SpecialInstructions = nil
	}
	if req.SpecialInstructions != nil && *req.SpecialInstructions == "" {
		req.SpecialInstructions = nil
	}
	// Prepare params with type conversions
	var menuIDStr string
	if req.MenuItemID != nil {
		menuIDStr = *req.MenuItemID
	}

	var specialInstructionsStr string
	if req.SpecialInstructions != nil {
		specialInstructionsStr = *req.SpecialInstructions
	}

	params := gen.CreateOrderItemParams{
		OrderID:             services.UUIDToPGType(menuItemID),
		MenuItemID:          services.StringToPGUUID(menuIDStr), // Handles empty string as NULL
		MenuItemName:        req.MenuItemName,
		Quantity:            req.Quantity,
		UnitPrice:           services.FloatToPGNumeric(req.UnitPrice),
		TotalPrice:          services.FloatToPGNumeric(req.TotalPrice),
		SpecialInstructions: services.StringToPGText(specialInstructionsStr),
	}
	orderItem, err := db.Q.CreateOrderItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to create order item",
		})
		return
	}
	c.JSON(http.StatusCreated, services.ToOrderItemResponse(orderItem))
}
func ListOrderItemsByOrderHandler(c *gin.Context) {
	idString := c.Param("order_items_id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	orderID, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}
	orderItems, err := db.Q.ListOrderItemsByOrder(c, services.UUIDToPGType(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list order items",
		})
		return
	}
	res := make([]models.OrderItem,0,len(orderItems))
	for _, orderItem := range orderItems {
		res = append(res, services.ToOrderItemResponse(orderItem))
	}
	c.JSON(http.StatusOK, res)
}

func UpdateOrderItemByIDHandler(c *gin.Context){
	idString := c.Param("order_items_id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order item id is required",
		})
		return
	}
	orderItemID, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order item id",
		})
		return
	}
	var req models.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	if *req.Quantity <= 0 || req.Quantity == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "quantity must be greater than 0",
		})
		return
	}
	if *req.UnitPrice <= 0 || req.UnitPrice == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "unit price must be greater than 0",
		})
		return
	}
	if *req.TotalPrice <= 0 || req.TotalPrice == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "total price must be greater than 0",
		})
		return
	}
	if *req.MenuItemName == "" || req.MenuItemName == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu item name is required",
		})
		return
	}
	if *req.SpecialInstructions == "" || req.SpecialInstructions == nil {
		req.SpecialInstructions = nil
	}
	// Prepare params with type conversions
	var menuIDStr string
	if req.MenuItemID != nil {
		menuIDStr = *req.MenuItemID
	}

	var specialInstructionsStr string
	if req.SpecialInstructions != nil {
		specialInstructionsStr = *req.SpecialInstructions
	}

	params := gen.UpdateOrderItemParams{
		ID:                services.UUIDToPGType(orderItemID),
		MenuItemID:        services.StringToPGUUID(menuIDStr), // Handles empty string as NULL
		MenuItemName:      *req.MenuItemName,
		Quantity:          *req.Quantity,
		UnitPrice:         services.FloatToPGNumeric(*req.UnitPrice),
		TotalPrice:        services.FloatToPGNumeric(*req.TotalPrice),
		SpecialInstructions: services.StringToPGText(specialInstructionsStr),
	}
	orderItem, err := db.Q.UpdateOrderItem(c, params)
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "order item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to update order item",
		})
		return
	}
	c.JSON(http.StatusOK, services.ToOrderItemResponse(orderItem))
}
func DeleteOrderItemHandler(c *gin.Context){
	idString := c.Param("order_items_id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order item id is required",
		})
		return
	}
	orderItemID, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order item id",
		})
		return
	}
	err = db.Q.DeleteOrderItem(c, services.UUIDToPGType(orderItemID))
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "order item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete order item",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "order item deleted successfully",
	})
}