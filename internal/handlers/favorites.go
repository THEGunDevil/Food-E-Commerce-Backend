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

func CreateFavoritesHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	userUUID, _ := uuid.Parse(userID)
	itemUUID, _ := uuid.Parse(req.MenuItemID.String())
	fav, err := db.Q.CreateFavorite(c, gen.CreateFavoriteParams{
		UserID:     services.UUIDToPGType(userUUID),
		MenuItemID: services.UUIDToPGType(itemUUID),
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Already favorite or invalid data",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Added favorite successfully!",
		Data:    services.ToFavResponse(fav),
	})
}

func DeleteFavoriteHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	menuIDStr := c.Param("menu_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid user id or user id required",
		})
	}
	menuUUID, err := uuid.Parse(menuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid menu id",
		})
	}
	err = db.Q.DeleteFavorite(c, gen.DeleteFavoriteParams{
		UserID:     services.UUIDToPGType(userUUID),
		MenuItemID: services.UUIDToPGType(menuUUID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no rows found to delete",
			})
		}
	}
    c.JSON(http.StatusOK,models.APIResponse{
        Success: true,
        Message: "One favorite item has been deleted successfully!",
    })
}

func ListFavoritesHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.CreateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	userUUID, _ := uuid.Parse(userID)
	fav, err := db.Q.ListFavoritesByUser(c, services.UUIDToPGType(userUUID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no favorites found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch favorites",
		})
		return
	}
	res := make([]models.FavoriteResponse, 0, len(fav))
	for _, f := range fav {
		res = append(res, services.ToFavListResponse(f))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "returned favorite list successfully!",
		Data:    res,
	})
}
