package handlers

import (
	"errors"
	"fmt"
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

func CreateCategoryHandler(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "name is required",
		})
		return
	}
	if req.Slug == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "slug is required",
		})
		return
	}
	if req.Description == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "description is required",
		})
		return
	}
	params := gen.CreateCategoryParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: services.StringToPGText(req.Description),
	}
	if req.DisplayOrder != nil {
		params.DisplayOrder = pgtype.Int4{Int32: *req.DisplayOrder, Valid: true}
	}
	if req.IsActive != nil {
		params.IsActive = services.PgTypeBool(req.IsActive)
	}

	if req.Image == nil {
		imageURL, publicID, err := services.HandleImageUpload(c, req.Image, "categories")
		if err != nil {
			return // uploadImage already sent error response
		}
		params.CatImageUrl = services.StringToPGText(imageURL)
		params.CatImagePublicID = services.StringToPGText(publicID)
	}
	cat, err := db.Q.CreateCategory(c, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category created successfully!",
		Data:    services.ToCategoryResponse(cat),
	})
}
func UpdateCategoryHandler(c *gin.Context) {
	idStr := c.Param("category_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "category id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid category id",
		})
		return
	}
	var req models.UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	params := gen.UpdateCategoryParams{
		ID: services.UUIDToPGType(id),
	}
	if req.CatImagePublicID != nil {
		params.CatImagePublicID = services.StringToPGText(*req.CatImagePublicID)
	}
	if req.Description != nil {
		params.Description = services.StringToPGText(*req.Description)
	}
	if req.DisplayOrder != nil {
		params.DisplayOrder = pgtype.Int4{Int32: *req.DisplayOrder, Valid: true}
	}
	if req.Slug != nil {
		params.Slug = *req.Slug
	}
	if req.Name != nil {
		params.Name = *req.Name
	}
	params.IsActive = services.PgTypeBool(req.IsActive)
	cat, err := db.Q.UpdateCategory(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "category not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category updated successfully!",
		Data:    services.ToCategoryResponse(cat),
	})
}
func UpdateCategoryDisplayOrderHandler(c *gin.Context) {
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
			Message: "invalid category id",
		})
		return
	}
	var req models.UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	params := gen.UpdateCategoryParams{
		ID: services.UUIDToPGType(catID),
	}

	if req.DisplayOrder != nil {
		params.DisplayOrder = pgtype.Int4{Int32: *req.DisplayOrder, Valid: true}
	} else {
		params.DisplayOrder = pgtype.Int4{Valid: false} // handle null if needed
	}

	cat, err := db.Q.UpdateCategory(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "category not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return // <- important
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category updated successfully!",
		Data:    services.ToCategoryResponse(cat),
	})
}
func GetCategoryByIDHandler(c *gin.Context){
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
			Message: "invalid category id",
		})
		return
	}
	cat, err := db.Q.GetCategoryByID(c, services.UUIDToPGType(catID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch category",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category returned successfully!",
		Data:    services.ToCategoryResponse(cat),
	})
}
func ListCategoriesHandler(c *gin.Context) {
	cats, err := db.Q.ListCategories(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no categories found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch categories",
		})
		return
	}
	res := make([]models.Category, 0, len(cats))
	for _, cat := range cats {
		res = append(res, services.ToCategoryResponse(cat))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Categories returned successfully!",
		Data:    res,
	})
}
func ListActiveCategoriesHandler(c *gin.Context) {
	cats, err := db.Q.ListActiveCategories(c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no active categories found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch active categories",
		})
		return
	}
	res := make([]models.Category, 0, len(cats))
	for _, cat := range cats {
		res = append(res, services.ToCategoryResponse(cat))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Active categories returned successfully!",
		Data:    res,
	})
}

func DeactivateCategoryHandler(c *gin.Context) {
	idStr := c.Param("category_id")
	catID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("failed to parse category id")
	}
	err = db.Q.DeactivateCategory(c, services.UUIDToPGType(catID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no category id found to delete",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to update category",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category deactivated successfully!",
	})
}

func DeleteCategoryHandler(c *gin.Context) {
	idStr := c.Param("category_id")
	catID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("failed to parse category id")
	}
	err = db.Q.DeleteCategory(c, services.UUIDToPGType(catID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no category id found to delete",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete category",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Deleted category successfully!",
	})
}
