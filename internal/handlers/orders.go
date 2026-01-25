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

func CreateOrderHandler(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		services.HandleValidationError(c, err)
		return
	}
	if req.UserID == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "user ID is required",
		})
		return
	}
	if req.DeliveryAddressID == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "delivery address ID is required",
		})
		return
	}
	if req.PaymentMethod == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "payment method is required",
		})
		return
	}
	if req.PaymentStatus == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "payment status is required",
		})
		return
	}
	if req.OrderStatus == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order status is required",
		})
		return
	}
	if req.TransactionID == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "transaction ID is required",
		})
		return
	}
	if req.CustomerName == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "customer name is required",
		})
		return
	}
	if req.CustomerPhone == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "customer phone is required",
		})
		return
	}
	if req.CustomerEmail == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "customer email is required",
		})
		return
	}
	if req.Subtotal == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "subtotal is required",
		})
		return
	}
	if req.DiscountAmount == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "discount amount is required",
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
	if req.VatAmount == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "VAT amount is required",
		})
		return
	}
	if req.TotalAmount == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "total amount is required",
		})
		return
	}
	if req.DeliveryPersonID == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "delivery person ID is required",
		})
		return
	}
	if req.EstimatedDelivery == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "estimated delivery is required",
		})
		return
	}
	if req.ActualDelivery == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "actual delivery is required",
		})
		return
	}
	if req.SpecialInstructions == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "special instructions is required",
		})
		return
	}
	if req.CancelledReason == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "cancelled reason is required",
		})
		return
	}
	if req.SpecialInstructions == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "special instructions is required",
		})
		return
	}
	if req.ActualDelivery == nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "actual delivery is required",
		})
		return
	}
	o, err := db.Q.CreateOrder(c, gen.CreateOrderParams{
		UserID:              services.UUIDToPGType(*req.UserID),
		DeliveryAddressID:   services.UUIDToPGType(*req.DeliveryAddressID),
		DeliveryType:        services.StringToPGText(req.DeliveryType),
		DeliveryAddress:     services.StringToPGText(req.DeliveryAddress),
		CustomerName:        req.CustomerName,
		CustomerPhone:       req.CustomerPhone,
		CustomerEmail:       services.StringToPGText(*req.CustomerEmail),
		Subtotal:            services.FloatToPGNumeric(req.Subtotal),
		DiscountAmount:      services.FloatToPGNumeric(*req.DiscountAmount),
		DeliveryFee:         services.FloatToPGNumeric(*req.DeliveryFee),
		VatAmount:           services.FloatToPGNumeric(*req.VatAmount),
		TotalAmount:         services.FloatToPGNumeric(req.TotalAmount),
		PaymentMethod:       services.StringToPGText(req.PaymentMethod),
		PaymentStatus:       services.StringToPGText(req.PaymentStatus),
		TransactionID:       services.StringToPGText(*req.TransactionID),
		OrderStatus:         services.StringToPGText(req.OrderStatus),
		DeliveryPersonID:    services.UUIDToPGType(*req.DeliveryPersonID),
		EstimatedDelivery:   pgtype.Timestamp{Time: *req.EstimatedDelivery, Valid: true},
		ActualDelivery:      pgtype.Timestamp{Time: *req.ActualDelivery, Valid: true},
		SpecialInstructions: services.StringToPGText(*req.SpecialInstructions),
		CancelledReason:     services.StringToPGText(*req.CancelledReason),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to create order",
		})
		return
	}
	
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "order created successfully",
		Data:    services.ToOrderResponse(o),
	})

}

func GetOrderByIDHandler(c *gin.Context) {
	idStr := c.Param("order_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}
	o, err := db.Q.GetOrderByID(c, services.UUIDToPGType(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to get order",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "order retrieved successfully",
		Data:    services.ToOrderResponse(o),
	})
}
func ListOrdersByUserHandler(c *gin.Context) {
	idStr := c.Param("user_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "user id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid user id",
		})
		return
	}
	r, err := db.Q.ListOrdersByUser(c, services.UUIDToPGType(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to list orders",
		})
		return
	}
	res := make([]models.Order, 0, len(r))
	for _, v := range r {
		res = append(res, services.ToOrderResponse(v))
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "orders listed successfully",
		Data:    res,
	})
}

func UpdateOrderByIDHandler(c *gin.Context) {
	idStr := c.Param("order_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}
	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "failed to bind request",
		})
		return
	}
	var params gen.UpdateOrderParams
	params.ID = services.UUIDToPGType(id)
	if req.DeliveryType != nil {
		params.DeliveryType = services.StringToPGText(*req.DeliveryType)
	}
	if req.DeliveryAddress != nil {
		params.DeliveryAddress = services.StringToPGText(*req.DeliveryAddress)
	}
	if req.CustomerName != nil {
		params.CustomerName = *req.CustomerName
	}
	if req.CustomerPhone != nil {
		params.CustomerPhone = *req.CustomerPhone
	}
	if req.CustomerEmail != nil {
		params.CustomerEmail = services.StringToPGText(*req.CustomerEmail)
	}
	if req.Subtotal != nil {
		params.Subtotal = services.FloatToPGNumeric(*req.Subtotal)
	}
	if req.DiscountAmount != nil {
		params.DiscountAmount = services.FloatToPGNumeric(*req.DiscountAmount)
	}
	if req.DeliveryFee != nil {
		params.DeliveryFee = services.FloatToPGNumeric(*req.DeliveryFee)
	}
	if req.VatAmount != nil {
		params.VatAmount = services.FloatToPGNumeric(*req.VatAmount)
	}
	if req.TotalAmount != nil {
		params.TotalAmount = services.FloatToPGNumeric(*req.TotalAmount)
	}
	if req.PaymentMethod != nil {
		params.PaymentMethod = services.StringToPGText(*req.PaymentMethod)
	}
	if req.PaymentStatus != nil {
		params.PaymentStatus = services.StringToPGText(*req.PaymentStatus)
	}
	if req.TransactionID != nil {
		params.TransactionID = services.StringToPGText(*req.TransactionID)
	}
	if req.OrderStatus != nil {
		params.OrderStatus = services.StringToPGText(*req.OrderStatus)
	}
	if req.DeliveryPersonID != nil {
		params.DeliveryPersonID = services.UUIDToPGType(*req.DeliveryPersonID)
	}
	if req.EstimatedDelivery != nil {
		params.EstimatedDelivery = pgtype.Timestamp{Time: *req.EstimatedDelivery, Valid: true}
	}
	if req.ActualDelivery != nil {
		params.ActualDelivery = pgtype.Timestamp{Time: *req.ActualDelivery, Valid: true}
	}
	if req.SpecialInstructions != nil {
		params.SpecialInstructions = services.StringToPGText(*req.SpecialInstructions)
	}
	if req.CancelledReason != nil {
		params.CancelledReason = services.StringToPGText(*req.CancelledReason)
	}
	if req.SpecialInstructions != nil {
		params.SpecialInstructions = services.StringToPGText(*req.SpecialInstructions)
	}
	if req.ActualDelivery != nil {
		params.ActualDelivery = pgtype.Timestamp{Time: *req.ActualDelivery, Valid: true}
	}
	o,err := db.Q.UpdateOrder(c,params)
	if err!= nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError,models.APIResponse{
			Success: false,
			Message: "failed to update order",
		})
		return
	}
	c.JSON(http.StatusOK,models.APIResponse{
		Success: true,
		Message: "order updated successfully",
		Data:    services.ToOrderResponse(o),
	})
}
func DeleteOrderByIDHandler(c *gin.Context) {
	idStr := c.Param("order_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "order id is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}
	err = db.Q.DeleteOrder(c, services.UUIDToPGType(id))
	if err != nil {
		if errors.Is(err,pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound,models.APIResponse{
				Success: false,
				Message: "order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to delete order",
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "order deleted successfully",
	})
}