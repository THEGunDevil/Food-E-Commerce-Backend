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
	"github.com/jackc/pgx/v5/pgtype"
)

func CreateCartItemsHandler(c *gin.Context) {
	var req models.CartItemsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	if req.MenuItemID == uuid.Nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu item id is required",
		})
		return
	}

	_, err := db.Q.GetCartItem(c, services.UUIDToPGType(req.MenuItemID))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Item does not exist → proceed to add
		} else {
			// Unexpected DB error
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "internal server error",
			})
			return
		}
	} else {
		// Item exists → return conflict
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "item already exists",
		})
		return
	}

	if req.Quantity < 1 || req.Quantity > 3 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "quantity must be between 1 and 3",
		})
		return
	}

	var specialInstructions string
	if req.SpecialInstructions != nil {
		specialInstructions = *req.SpecialInstructions
	}

	var userUUID pgtype.UUID
	var sessionUUID pgtype.UUID

	userIDStr := c.GetString("user_id")

	if userIDStr != "" {
		// Authenticated user
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "invalid user id",
			})
			return
		}

		userUUID = pgtype.UUID{
			Bytes: uid,
			Valid: true,
		}
		sessionUUID = pgtype.UUID{Valid: false}

	} else {
		sid := services.GetOrCreateSessionID(c)
		if sid == "" {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "failed to create session",
			})
			return
		}

		parsedSID, err := uuid.Parse(sid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "invalid session id",
			})
			return
		}

		userUUID = pgtype.UUID{Valid: false}
		sessionUUID = pgtype.UUID{
			Bytes: parsedSID,
			Valid: true,
		}
	}

	params := gen.AddCartItemParams{
		UserID:              userUUID,    // NULL or UUID
		SessionID:           sessionUUID, // NULL or UUID
		MenuItemID:          pgtype.UUID{Bytes: req.MenuItemID, Valid: true},
		Quantity:            int32(req.Quantity),
		SpecialInstructions: services.StringToPGText(specialInstructions),
	}

	res, err := db.Q.AddCartItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "cart item added successfully",
		Data:    services.ToCartItemResponse(res),
	})
}

func ListCartItemsHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "user_id is required"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "user_id is required"})
		return
	}
	cartItems, err := db.Q.ListCartItemsByUser(c, services.UUIDToPGType(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "failed to list cart items"})
		return
	}
	res := make([]models.CartItems, 0, len(cartItems))
	for _, cartItem := range cartItems {
		res = append(res, services.ToCartItemResponse(cartItem))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "cart items listed successfully",
		Data:    res,
	})
}
func UpdateCartItemHandler(c *gin.Context) {
	cartItemIDStr := c.Param("cart_item_id")
	if cartItemIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "cart item id is required"})
		return
	}
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "cart item id is required"})
		return
	}
	var req models.UpdateCartItemQuantityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quantity := req.Quantity
	if quantity == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "quantity is required"})
		return
	}
	specialInstructions := req.SpecialInstructions
	if specialInstructions == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "special instructions is required"})
		return
	}
	params := gen.UpdateCartItemParams{
		Quantity:            int32(*quantity),
		SpecialInstructions: services.StringToPGText(*specialInstructions),
		ID:                  services.UUIDToPGType(cartItemID),
	}
	res, err := db.Q.UpdateCartItem(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "failed to update cart item quantity"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "cart item quantity updated successfully",
		Data:    services.ToCartItemResponse(res),
	})
}
func RemoveCartItemHandler(c *gin.Context) {
	cartItemIDStr := c.Param("cart_item_id")
	if cartItemIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "cart item id is required"})
		return
	}
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "cart item id is required"})
		return
	}
	err = db.Q.RemoveCartItem(c, services.UUIDToPGType(cartItemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "cart item not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "failed to remove cart item"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "cart item removed successfully"})
}
func ClearCartHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "user id is required"})
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "user id is required"})
		return
	}
	err = db.Q.ClearCartByUser(c, services.UUIDToPGType(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "cart not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{Success: false, Message: "failed to clear cart"})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{Success: true, Message: "cart cleared successfully"})
}
