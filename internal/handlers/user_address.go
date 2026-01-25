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
)

func CreateUserAddressHandler(c *gin.Context) {
	var req models.CreateUserAddressRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	if req.AddressLine1 == ""{
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "address line 1 is required",
		})
		return
	}
	if req.AddressLine2 == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "address line 2 is required",
		})
		return
	}
	if req.Area == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "area is required",
		})
		return
	}
	if req.City == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "city is required",
		})
		return
	}
	if req.IsDefault == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "is default is required",
		})
		return
	}
	if req.Label == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "label is required",
		})
		return
	}
	if req.Latitude == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "latitude is required",
		})
		return
	}
	if req.Longitude == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "longitude is required",
		})
		return
	}
	if req.PostalCode == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "postal code is required",
		})
		return
	}
	ua, err := db.Q.CreateUserAddress(c, gen.CreateUserAddressParams{
		AddressLine1: req.AddressLine1,
		AddressLine2: services.StringToPGText(*req.AddressLine2),
		Area:         req.Area,
		City:         req.City,
		IsDefault:    req.IsDefault,
		Label:        req.Label,
		Latitude:     services.FloatToPGNumeric(*req.Latitude),
		Longitude:    services.FloatToPGNumeric(*req.Longitude),
		PostalCode:   services.StringToPGText(*req.PostalCode),
		UserID:       services.StringToPGUUID(req.UserID),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "failed to create user address",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User address created successfully!",
		Data:    services.ToUserAddressResponse(ua),
	})
}
func UpdateUserAddressByIDHandler(c *gin.Context) {
	idStr := c.Param("userid")
	if idStr == "" {
		fmt.Println("id required to update an address")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("failed to parse id")
	}
	var req models.UpdateUserAddressRequest
	if err := c.ShouldBind(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	params := gen.UpdateUserAddressParams{
		ID: services.UUIDToPGType(id),
	}
	if req.AddressLine1 != nil {
		params.AddressLine1 = services.StringToPGText(*req.AddressLine1)
	}
	if req.AddressLine2 != nil {
		params.AddressLine2 = services.StringToPGText(*req.AddressLine2)
	}
	if req.Area != nil {
		params.Area = services.StringToPGText(*req.Area)
	}
	if req.City != nil {
		params.City = services.StringToPGText(*req.City)
	}
	params.IsDefault = services.PgTypeBool(req.IsDefault)
	if req.Label != nil {
		params.Label = services.StringToPGText(*req.Label)
	}
	if req.Latitude != nil {
		params.Latitude = services.FloatToPGNumeric(req.Latitude)
	}
	if req.Longitude != nil {
		params.Longitude = services.FloatToPGNumeric(req.Longitude)
	}
	if req.PostalCode != nil {
		params.PostalCode = services.StringToPGText(*req.PostalCode)
	}
	ua, err := db.Q.UpdateUserAddress(c, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no address found matching this id",
			})
			return
		}
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "failed to update user address",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User address updated successfully!",
		Data:    services.ToUserAddressResponse(ua),
	})
}
func GetUserAddressByIDHandler(c *gin.Context) {
	idStr := c.GetString("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println("failed to parse user id")
	}
	u, err := db.Q.GetUserAddressByID(c, services.UUIDToPGType(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no address found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch user address",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User address fetch successfully!",
		Data:    services.ToUserAddressResponse(u),
	})
}
func DeleteUserAddressByIDHandler(c *gin.Context){
	idStr:= c.Param("userid")
	if idStr == "" {
		fmt.Println("id required to delete an address")
		return
	}
	id,err := uuid.Parse(idStr)
	if err!= nil {
		fmt.Println("failed to parse id")
	}
	err=db.Q.DeleteUserAddress(c, services.UUIDToPGType(id))
	if err!=nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "no address found matching this id",
			})
			return
		}
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to delete user address, err:",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Address deleted successfully!",
	})
}