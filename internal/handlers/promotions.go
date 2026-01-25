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

func CreatePromotionHandler(c *gin.Context) {
	var req models.CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}

	if req.Title != "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "title is required",
		})
		return
	}
	if req.Description == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "description is required",
		})
		return
	}
	if req.DiscountType == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "discount type is required",
		})
		return
	}
	if req.DiscountValue == nil || *req.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "discount value is required",
		})
		return
	}
	if req.MinOrderAmount == nil || *req.MinOrderAmount <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "min order amount is required",
		})
		return
	}
	if req.ValidFrom.IsZero() {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "valid from is required",
		})
		return
	}

	if req.ValidUntil.IsZero() {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "valid until is required",
		})
		return
	}
	if req.MaxUses == nil || *req.MaxUses <= 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "max uses is required",
		})
		return
	}

	params := gen.CreatePromotionParams{
		Title:         req.Title,
		Description:   services.StringToPGText(*req.Description),
		DiscountType:  services.StringToPGText(req.DiscountType),
		DiscountValue: services.FloatToPGNumeric(req.DiscountValue),
		ValidFrom:     services.TimeToTimestamp(req.ValidFrom),
		ValidUntil:    services.TimeToTimestamp(req.ValidUntil),
		MaxUses:       pgtype.Int4{Int32: *req.MaxUses, Valid: true},
		IsActive:      services.PgTypeBool(req.IsActive),
	}

	p, err := db.Q.CreatePromotion(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to create promotion",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "promotion created successfully",
		Data:    services.ToPromotionResponse(p),
	})
}
func ListPromotionsHandler(c *gin.Context) {
	promotions, err := db.Q.ListPromotions(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list promotions",
		})
		return
	}
	res:=make([]models.Promotion, 0,len(promotions))
	for _,promotion:=range promotions {
		res=append(res,services.ToPromotionResponse(promotion))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "promotions listed successfully",
		Data:    res,
	})
}

func ListActivePromotionsHandler(c *gin.Context){
	promotions,err:=db.Q.ListActivePromotions(c)
	if err!=nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusOK,models.APIResponse{
				Success: true,
				Message: "no active promotions found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to list active promotions",
		})
		return
	}
	res:=make([]models.Promotion, 0,len(promotions))
	for _,promotion:=range promotions {
		res=append(res,services.ToPromotionResponse(promotion))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "active promotions listed successfully",
		Data:    res,
	})
}
func UpdatePromotionHandler(c *gin.Context){
	idStr:=c.Param("promotion_id")
	if idStr=="" {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "promotion id is required",
		})
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	var req models.UpdatePromotionRequest
	if err:=c.ShouldBindJSON(&req);err!=nil {
		services.HandleValidationError(c,err)
		return
	}
	params:=gen.UpdatePromotionParams{
		ID: services.UUIDToPGType(id),
	}
	if req.Title!=nil {
		params.Title=*req.Title
	}
	if req.Description!=nil {
		params.Description=services.StringToPGText(*req.Description)
	}
	if req.DiscountType!=nil {
		params.DiscountType=services.StringToPGText(*req.DiscountType)
	}
	if req.DiscountValue!=nil {
		params.DiscountValue=services.FloatToPGNumeric(*req.DiscountValue)
	}
	if req.MinOrderAmount!=nil {
		params.MinOrderAmount=services.FloatToPGNumeric(*req.MinOrderAmount)
	}
	if req.ValidFrom.IsZero() {
		params.ValidFrom=services.TimeToTimestamp(*req.ValidFrom)
	}
	if req.ValidUntil.IsZero() {
		params.ValidUntil=services.TimeToTimestamp(*req.ValidUntil)
	}
	if req.MaxUses!=nil {
		params.MaxUses=pgtype.Int4{Int32: *req.MaxUses, Valid: true}
	}
	if req.IsActive!=nil {
		params.IsActive=services.PgTypeBool(req.IsActive)
	}
	p,err:=db.Q.UpdatePromotion(c,params)
	if err!=nil {
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to update promotion",
		})
		return
	}
	c.JSON(http.StatusOK,models.APIResponse{
		Success: true,
		Message: "promotion updated successfully",
		Data:    services.ToPromotionResponse(p),
	})
}
func DeletePromotionHandler(c *gin.Context){
	idStr:=c.Param("promotion_id")
	if idStr=="" {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "promotion id is required",
		})
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	err=db.Q.DeletePromotion(c,services.UUIDToPGType(id))
	if err!=nil {
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to delete promotion",
		})
		return
	}
	c.JSON(http.StatusOK,models.APIResponse{
		Success: true,
		Message: "promotion deleted successfully",
	})
}
func IncrementPromotionUsageHandler(c *gin.Context){
	idStr:=c.Param("promotion_id")
	if idStr=="" {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "promotion id is required",
		})
		return
	}
	id,err:=uuid.Parse(idStr)
	if err!=nil {
		c.JSON(http.StatusBadRequest,models.APIResponse{
			Success: false,
			Message: "invalid id",
		})
		return
	}
	p,err:=db.Q.IncrementPromotionUsage(c,services.UUIDToPGType(id))
	if err!=nil {
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to increment promotion usage",
		})
		return
	}
	c.JSON(http.StatusOK,models.APIResponse{
		Success: true,
		Message: "promotion usage incremented successfully",
		Data:    services.ToPromotionResponse(p),
	})
}