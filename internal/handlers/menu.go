package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateMenuItemHandler handles creating a new menu item
func CreateMenuItemHandler(c *gin.Context) {
	var req models.CreateMenuItemMultipartRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	// 1. Build Item Params (passing empty slice since images aren't uploaded yet)
	params, _, err := services.BuildCreateMenuItemParams(req.CategoryID, &req, []models.MenuItemImage{})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 2. Insert Menu Item into DB
	menuItem, err := db.Q.CreateMenuItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create menu item",
		})
		return
	}

	// 3. Handle image upload
	var dbImages []gen.MenuItemImage // Use the generator's type for the response builder
	if req.Image != nil {
		imageURL, publicID, err := services.HandleImageUpload(c, req.Image, "menu_items")
		if err != nil {
			// Failure to upload image shouldn't necessarily crash the whole request,
			// but since it's required in your struct, we return.
			return
		}

		imgParams := gen.CreateMenuItemImageParams{
			MenuItemID:    menuItem.ID,
			ImageUrl:      imageURL,
			ImagePublicID: services.StringToPGText(publicID),
			IsPrimary:     pgtype.Bool{Bool: true, Valid: true}, // Fixed boolean logic
			DisplayOrder:  pgtype.Int4{Int32: 0, Valid: true},
		}

		// Save to DB
		err = db.Q.CreateMenuItemImage(c, imgParams)
		if err == nil {
			// Fetch the images associated with this item to include in response
			dbImages, _ = db.Q.ListMenuItemImages(c, menuItem.ID)
		}
	}

	// 4. Map to final response
	// ToMenuResponse usually takes the DB model and the DB image list
	resp := services.ToMenuResponse(menuItem, dbImages)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Menu item created successfully",
		Data:    resp,
	})
}

// UpdateMenuByMenuIDHandler handles updating an existing menu item
func UpdateMenuByMenuIDHandler(c *gin.Context) {
	id := c.Param("menu_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu id required",
		})
		return
	}

	var req models.UpdateMenuItemMultipartRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	categoryID := id
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	params, err := services.BuildUpdateMenuItemParams(categoryID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Ensure the ID for the update query is set
	menuUUID := services.StringToPGUUID(id) // Helper assuming string->uuid
	params.ID = menuUUID

	// Logic for Image Replacement/Deletion
	shouldDeleteOldImage := (req.Image != nil || (req.RemoveImage != nil && *req.RemoveImage))

	if shouldDeleteOldImage && req.ImagePublicID != nil && *req.ImagePublicID != "" {
		services.DeleteImageFromCloudinary(*req.ImagePublicID)

	}

	// Handle New Image Upload
	if req.Image != nil {
		imageURL, publicID, err := services.HandleImageUpload(c, req.Image, "menu_items")
		if err != nil {
			return
		}

		imgParams := gen.CreateMenuItemImageParams{
			MenuItemID:    menuUUID,
			ImageUrl:      imageURL,
			ImagePublicID: services.StringToPGText(publicID),
			IsPrimary:     pgtype.Bool{Bool: true, Valid: true},
			DisplayOrder:  pgtype.Int4{Int32: 0, Valid: true},
		}
		db.Q.CreateMenuItemImage(c, imgParams)
	}

	// Update in DB
	res, err := db.Q.UpdateMenuItem(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update menu item",
		})
		return
	}

	// Fetch images to return complete response
	dbImages, _ := db.Q.ListMenuItemImages(c, res.ID)
	// _ = services.ToMenuItemImages(dbImages)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Menu item updated successfully",
		Data:    services.ToMenuResponse(res, dbImages),
	})
}

func DeleteMenuByMenuIDHandler(c *gin.Context) {
	id := c.Param("menu_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu id required",
		})
		return
	}

	// Validate UUID format
	pgUUID := services.StringToPGUUID(id)
	if !pgUUID.Valid {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid menu id format",
		})
		return
	}

	err := db.Q.DeleteMenuItem(c, pgUUID)
	if err != nil {
		// Handle not found error
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "menu item not found",
			})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete menu item",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "menu item deleted successfully",
	})
}

func ListMenusHandler(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	menus, err := db.Q.ListMenuItems(c, gen.ListMenuItemsParams{Limit: int32(limit),
		Offset: int32(offset)})

	if err != nil {
		// Handle not found error
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "menu items not found",
			})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	menuCount, err := db.Q.CountMenuItems(c)
	if err != nil {
		c.JSON(http.StatusNoContent, models.APIResponse{
			Success: false,
			Message: "menu is empty",
			Error:   err.Error(),
		})
		return
	}
	totalPages := int(math.Ceil(float64(menuCount) / float64(limit)))

	response := make([]models.MenuItem, 0, len(menus))
	for _, m := range menus {
		dbImages, _ := db.Q.GetMenuItemImagesByMenuItemID(c, m.ID)
		response = append(response, services.ToMenuListResponseWithCategoryName(m, dbImages))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:    true,
		Page:       page,
		Message: "Menus successfully returned!",
		Limit:      limit,
		TotalPages: totalPages,
		Data:       response,
	})
}
func ListMenusByCategoryHandler(c *gin.Context) {
	idStr := c.Param("category_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "category id is required",
		})
		return
	}
	catID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	menus, err := db.Q.ListMenuItemsByCategory(c, gen.ListMenuItemsByCategoryParams{
		CategoryID: services.UUIDToPGType(catID),
		Limit:      int32(limit),
		Offset:     int32(offset)})

	if err != nil {
		// Handle not found error
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "menu items not found",
			})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	menuCount, err := db.Q.CountMenuItemsByCategory(c, services.UUIDToPGType(catID))
	if err != nil {
		c.JSON(http.StatusNoContent, models.APIResponse{
			Success: false,
			Message: "menu is empty",
			Error:   err.Error(),
		})
		return
	}
	totalPages := int(math.Ceil(float64(menuCount) / float64(limit)))

	response := make([]models.MenuItem, 0, len(menus))
	for _, m := range menus {
		dbImages, _ := db.Q.GetMenuItemImagesByMenuItemID(c, m.ID)
		response = append(response, services.ToMenuListResponseWithCategoryName(gen.ListMenuItemsRow(m), dbImages))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:    true,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       response,
	})
}
func GetMenuByMenuIDHandler(c *gin.Context) {
	id := c.Param("menu_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu id required",
		})
		return
	}

	// Validate UUID format
	pgUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid menu id format",
		})
		return
	}

	menu, err := db.Q.GetMenuItemByID(c, services.UUIDToPGType(pgUUID))
	if err != nil {
		// Handle not found error
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "menu item not found",
			})
			return
		}

		// Handle other database errors
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	dbImages, _ := db.Q.ListMenuItemImages(c, services.UUIDToPGType(pgUUID))
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Menu item retrieved successfully!!",
		Data:    services.ToMenuListResponseWithCategoryName(gen.ListMenuItemsRow(menu), dbImages),
	})
}

