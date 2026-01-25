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
)

func CreateReviewHandler(c *gin.Context) {
	idStr := c.Param("menu_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu id is required",
		})
		return
	}
	menuItemID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	var req models.CreateReviewRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "rating must be between 1 and 5",
		})
		return
	}
	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "user id is required",
		})
		return
	}
	if req.OrderID == uuid.Nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	if req.Comment == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "comment is required",
		})
		return
	}
	if len(req.Images) > 5 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "max 5 images allowed",
		})
		return
	}
	params := gen.CreateReviewParams{
		MenuItemID: services.UUIDToPGType(menuItemID),
		UserID:     services.UUIDToPGType(req.UserID),
		Rating:     *req.Rating,
		Comment:    services.StringToPGText(*req.Comment),
	}
	review, err := db.Q.CreateReview(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to create review",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "review created successfully",
		Data:    services.ToReviewResponse(review),
	})
}
func GetReviewByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "id is required",
		})
		return
	}
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}

	review, err := db.Q.GetReviewByID(c, services.UUIDToPGType(reviewID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to get review",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "review retrieved successfully",
		Data:    services.ToReviewResponse(review),
	})
}
func GetMenuItemReviewsByMenuItemIDHandler(c *gin.Context) {
	idStr := c.Param("menu_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "menu id is required",
		})
		return
	}
	menuItemID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid menu id",
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

	reviews, err := db.Q.ListMenuItemReviewsByMenuItemID(c, gen.ListMenuItemReviewsByMenuItemIDParams{
		MenuItemID: services.UUIDToPGType(menuItemID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "reviews not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to get reviews",
		})
		return
	}
	totalReviews, err := db.Q.CountMenuItemReviewsByMenuItemID(c, services.UUIDToPGType(menuItemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to get reviews",
		})
		return
	}
	totalPages := int(math.Ceil(float64(totalReviews) / float64(limit)))
	response := make([]models.Review, 0, len(reviews))
	for _, review := range reviews {
		response = append(response, services.ToReviewResponse(review))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success:    true,
		Message:    "review retrieved successfully",
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       response,
	})
}

func ListReviewsByUserHandler(c *gin.Context) {
	idStr := c.Param("user_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "user id is required",
		})
		return
	}
	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	reviews, err := db.Q.ListReviewsByUser(c, services.UUIDToPGType(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list reviews",
		})
		return
	}
	res := make([]models.Review, 0, len(reviews))
	for _, review := range reviews {
		res = append(res, services.ToReviewResponse(review))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "reviews listed successfully",
		Data:    res,
	})
}

// func UpdateReviewByIDHandler(c *gin.Context) {
// 	idStr := c.Param("review_id")
// 	if idStr == "" {
// 		c.JSON(http.StatusBadRequest, models.APIResponse{
// 			Success: false,
// 			Message: "review id is required",
// 		})
// 		return
// 	}
// 	reviewID, err := uuid.Parse(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, models.APIResponse{
// 			Success: false,
// 			Message: "invalid id",
// 		})
// 		return
// 	}
// 	var req models.UpdateReviewRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		services.HandleValidationError(c, err)
// 		return
// 	}
// 	params := gen.UpdateReviewParams{
// 		ID: services.UUIDToPGType(reviewID),
// 	}
// 	if req.Rating != nil {
// 		params.Rating = *req.Rating
// 	}
// 	if req.Comment != nil {
// 		params.Comment = services.StringToPGText(*req.Comment)
// 	}
// 	if len(req.Images) > 5 {
// 		c.JSON(http.StatusBadRequest, models.APIResponse{
// 			Success: false,
// 			Message: "max 5 images allowed",
// 		})
// 		return
// 	}
// 	params.Images = req.Images
// 	review, err := db.Q.UpdateReview(c, params)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			c.JSON(http.StatusNotFound, models.APIResponse{
// 				Success: false,
// 				Message: "review not found",
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, models.APIResponse{
// 			Success: false,
// 			Message: "failed to update review",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, models.APIResponse{
// 		Success: true,
// 		Message: "review updated successfully",
// 		Data:    services.ToReviewResponse(review),
// 	})
// }

func DeleteReviewByIDHandler(c *gin.Context) {
	idStr := c.Param("review_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "review id is required",
		})
		return
	}
	reviewID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	err = db.Q.DeleteReview(c, services.UUIDToPGType(reviewID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "review not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete review",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "review deleted successfully",
	})
}
