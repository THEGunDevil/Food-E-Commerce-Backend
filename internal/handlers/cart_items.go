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

func CreateCartItemsHandler(c *gin.Context) {
	var req models.CartItemsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	// 1️⃣ Validate menu item
	if req.MenuItemID == uuid.Nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu item id is required",
		})
		return
	}

	// 2️⃣ Get cart_id from middleware
	cartIDValue, exists := c.Get("cart_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "cart not initialized",
		})
		return
	}
	cartID, ok := cartIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "invalid cart id type",
		})
		return
	}

	// 3️⃣ Quantity validation
	if req.Quantity < 1 || req.Quantity > 3 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "quantity must be between 1 and 3",
		})
		return
	}

	// 4️⃣ Check if item already exists
	_, err := db.Q.GetCartItemByCartAndMenuItem(
		c,
		gen.GetCartItemByCartAndMenuItemParams{
			CartID:     services.UUIDToPGType(cartID),
			MenuItemID: services.UUIDToPGType(req.MenuItemID),
		},
	)

	if err == nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "item already exists in cart",
		})
		return
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to check cart item",
			Error:   err.Error(),
		})
		return
	}

	// 5️⃣ Prepare insert params
	var instructions string
	if req.SpecialInstructions != nil {
		instructions = *req.SpecialInstructions
	}

	params := gen.AddCartItemParams{
		CartID:              services.UUIDToPGType(cartID),
		MenuItemID:          services.UUIDToPGType(req.MenuItemID),
		Quantity:            int32(req.Quantity),
		SpecialInstructions: services.StringToPGText(instructions),
	}

	// 6️⃣ Insert cart item
	_, err = db.Q.AddCartItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to add item to cart",
			Error:   err.Error(),
		})
		return
	}

	// 7️⃣ Success
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "cart item added successfully",
	})
}

func ListCartItemsHandler(c *gin.Context) {
	cartIDValue, exists := c.Get("cart_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "cart not found",
		})
		return
	}

	cartID, ok := cartIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "invalid cart id",
		})
		return
	}

	cartItems, err := db.Q.ListCartItemsByCart(
		c,
		services.UUIDToPGType(cartID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list cart items",
		})
		return
	}

	res := make([]models.CartItemResponse, 0, len(cartItems))

	for _, cartItem := range cartItems {
		menuItem, err := db.Q.GetMenuItemByID(c, cartItem.MenuItemID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, models.APIResponse{
					Success: false,
					Message: "menu item not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "error fetching menu item",
			})
			return
		}

		menuImages, err := db.Q.ListMenuItemImages(c, menuItem.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "error fetching menu images",
			})
			return
		}

		res = append(res,
			services.ToCartItemResponse(menuItem, cartItem, menuImages),
		)
	}

	subTotal, err := db.Q.GetSubTotal(c, services.UUIDToPGType(cartID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "there was an error calculating subtotal",
			Error:   err.Error(),
		})
	}
	subTotalFloat, _ := subTotal.Float64Value()
	c.JSON(http.StatusOK, models.APIResponse{
		Success:    true,
		Message:    "Cart items listed successfully",
		TotalItems: len(res),
		SubTotal:   subTotalFloat.Float64,
		Data:       res,
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
	_, err = db.Q.UpdateCartItem(c, params)
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
		// Data:    services.ToCartItemResponse(res),
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
	// 1️⃣ Get cart_id from context (set by SessionMiddleware)
	cartIDValue, exists := c.Get("cart_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "cart not found",
		})
		return
	}

	cartID, ok := cartIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "invalid cart id",
		})
		return
	}

	// 2️⃣ Delete all items from this cart
	err := db.Q.ClearCart(c, services.UUIDToPGType(cartID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to clear cart",
		})
		return
	}

	// 3️⃣ Respond
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "cart cleared successfully",
	})
}
