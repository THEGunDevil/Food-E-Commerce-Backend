package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID
	OrderNumber         string
	UserID              *uuid.UUID
	DeliveryAddressID   *uuid.UUID
	DeliveryType        string
	DeliveryAddress     string
	CustomerName        string
	CustomerPhone       string
	CustomerEmail       *string
	Subtotal            float64
	DiscountAmount      float64
	DeliveryFee         float64
	VATAmount           float64
	TotalAmount         float64
	PaymentMethod       string
	PaymentStatus       string
	TransactionID       *string
	OrderStatus         string
	DeliveryPersonID    *uuid.UUID
	EstimatedDelivery   *time.Time
	ActualDelivery      *time.Time
	SpecialInstructions *string
	CancelledReason     *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateOrderRequest struct {
	UserID              *uuid.UUID `json:"user_id"`             // optional
	DeliveryAddressID   *uuid.UUID `json:"delivery_address_id"` // optional
	DeliveryType        string     `json:"delivery_type" binding:"required,oneof=delivery pickup"`
	DeliveryAddress     string     `json:"delivery_address" binding:"required"`
	CustomerName        string     `json:"customer_name" binding:"required"`
	CustomerPhone       string     `json:"customer_phone" binding:"required"`
	CustomerEmail       *string    `json:"customer_email"`
	Subtotal            float64    `json:"subtotal" binding:"required,gt=0"`
	DiscountAmount      *float64   `json:"discount_amount"`
	DeliveryFee         *float64   `json:"delivery_fee"`
	VatAmount           *float64   `json:"vat_amount"`
	TotalAmount         float64    `json:"total_amount" binding:"required,gt=0"`
	PaymentMethod       string     `json:"payment_method" binding:"required,oneof=cod bkash nagad card rocket"`
	PaymentStatus       string     `json:"payment_status" binding:"required,oneof=pending paid failed refunded"`
	TransactionID       *string    `json:"transaction_id"`
	OrderStatus         string     `json:"order_status" binding:"required,oneof=pending confirmed preparing ready on_the_way delivered cancelled failed"`
	DeliveryPersonID    *uuid.UUID `json:"delivery_person_id"`
	EstimatedDelivery   *time.Time `json:"estimated_delivery"`
	ActualDelivery      *time.Time `json:"actual_delivery"`
	SpecialInstructions *string    `json:"special_instructions"`
	CancelledReason     *string    `json:"cancelled_reason"`
}
type UpdateOrderRequest struct {
	DeliveryType        *string    `json:"delivery_type" binding:"omitempty,oneof=delivery pickup"`
	DeliveryAddress     *string    `json:"delivery_address"`
	CustomerName        *string    `json:"customer_name"`
	CustomerPhone       *string    `json:"customer_phone"`
	CustomerEmail       *string    `json:"customer_email"`
	Subtotal            *float64   `json:"subtotal" binding:"omitempty,gt=0"`
	DiscountAmount      *float64   `json:"discount_amount"`
	DeliveryFee         *float64   `json:"delivery_fee"`
	VatAmount           *float64   `json:"vat_amount"`
	TotalAmount         *float64   `json:"total_amount" binding:"omitempty,gt=0"`
	PaymentMethod       *string    `json:"payment_method" binding:"omitempty,oneof=cod bkash nagad card rocket"`
	PaymentStatus       *string    `json:"payment_status" binding:"omitempty,oneof=pending paid failed refunded"`
	TransactionID       *string    `json:"transaction_id"`
	OrderStatus         *string    `json:"order_status" binding:"omitempty,oneof=pending confirmed preparing ready on_the_way delivered cancelled failed"`
	DeliveryPersonID    *uuid.UUID `json:"delivery_person_id"`
	EstimatedDelivery   *time.Time `json:"estimated_delivery"`
	ActualDelivery      *time.Time `json:"actual_delivery"`
	SpecialInstructions *string    `json:"special_instructions"`
	CancelledReason     *string    `json:"cancelled_reason"`
}
