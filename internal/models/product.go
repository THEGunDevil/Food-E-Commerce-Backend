package models

import (
    "mime/multipart"
    "time"

    "github.com/google/uuid"
)

type MenuItemImage struct {
    ID            uuid.UUID `json:"id"`
    ImageUrl      string    `json:"image_url"`
    ImagePublicID string    `json:"image_public_id"` // Added field
    IsPrimary     bool      `json:"is_primary"`
    DisplayOrder  int     `json:"display_order"`
}

type MenuItem struct {
    ID            uuid.UUID       `json:"id"`
    CategoryName  string          `json:"category_name"`
    CategoryID    uuid.UUID       `json:"category_id"`
    Name          string          `json:"name"`
    Slug          string          `json:"slug"`
    Description   string          `json:"description"`
    Price         float64         `json:"price"`
    DiscountPrice *float64         `json:"discount_price"`
    Images        []MenuItemImage `json:"images"`
    Ingredients   []string        `json:"ingredients"`
    Tags          []string        `json:"tags"`
    PrepTime      int           `json:"prep_time"`
    SpicyLevel    int            `json:"spicy_level"`
    IsVegetarian  bool            `json:"is_vegetarian"`
    IsSpecial     bool            `json:"is_special"`
    IsAvailable   bool            `json:"is_available"`
    StockQuantity int           `json:"stock_quantity"`
    MinStockAlert int           `json:"min_stock_alert"`
    TotalOrders   int           `json:"total_orders"`
    AverageRating float64         `json:"average_rating"`
    DisplayOrder  int           `json:"display_order"`
    CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
}

type CreateMenuItemMultipartRequest struct {
    CategoryID    string                `form:"category_id" validate:"required,uuid4"`
    Name          string                `form:"name" validate:"required,min=2,max=200"`
    Slug          string                `form:"slug" validate:"required,min=2,max=200,slug"`
    Description   string                `form:"description" validate:"max=1000"`
    Price         string                `form:"price" validate:"required,gt=0"`
    DiscountPrice *string               `form:"discount_price" validate:"omitempty,gte=0"`
    Ingredients   []string              `form:"ingredients[]"`
    Tags          []string              `form:"tags[]"`
    PrepTime      *int                `form:"prepTime" validate:"omitempty,gte=0,lte=1440"`
    SpicyLevel    *int                 `form:"spicy_level" validate:"omitempty,gte=0,lte=3"`
    IsVegetarian  *bool                 `form:"is_vegetarian"`
    IsSpecial     *bool                 `form:"is_special"`
    IsAvailable   *bool                 `form:"is_available"`
    StockQuantity *int                `form:"stock_quantity" validate:"omitempty,gte=-1"`
    MinStockAlert *int                `form:"min_stock_alert" validate:"omitempty,gte=0"`
    DisplayOrder  *int                `form:"display_order" validate:"omitempty,gte=0"`
    Image         *multipart.FileHeader `form:"image" validate:"required"`
}

type UpdateMenuItemMultipartRequest struct {
    MenuItemID    string                `form:"-"` // From URL parameter
    CategoryID    *string               `form:"category_id" validate:"omitempty,uuid4"`
    Name          *string               `form:"name" validate:"omitempty,min=2,max=200"`
    Slug          *string               `form:"slug" validate:"omitempty,min=2,max=200,slug"`
    Description   *string               `form:"description" validate:"omitempty,max=1000"`
    Price         *string               `form:"price" validate:"omitempty,gt=0"`
    DiscountPrice *string               `form:"discount_price" validate:"omitempty,gte=0"`
    Ingredients   []string              `form:"ingredients[]" validate:"omitempty,max=50,dive,min=1,max=100"`
    Tags          []string              `form:"tags[]" validate:"omitempty,max=20,dive,min=1,max=50"`
    PrepTime      *int                `form:"prep_time" validate:"omitempty,gte=0,lte=1440"`
    SpicyLevel    *int8                 `form:"spicy_level" validate:"omitempty,gte=0,lte=3"`
    IsVegetarian  *bool                 `form:"is_vegetarian"`
    IsSpecial     *bool                 `form:"is_special"`
    IsAvailable   *bool                 `form:"is_available"`
    StockQuantity *int                `form:"stock_quantity" validate:"omitempty,gte=-1"`
    MinStockAlert *int                `form:"min_stock_alert" validate:"omitempty,gte=0"`
    DisplayOrder  *int                `form:"display_order" validate:"omitempty,gte=0"`
    
    // Image Handling
    Image         *multipart.FileHeader `form:"image" validate:"omitempty"`           // New image file
    ImagePublicID *string               `form:"image_public_id" validate:"omitempty"` // ID of the OLD image to delete
    RemoveImage   *bool                 `form:"remove_image" validate:"omitempty"`     // Explicit flag to remove old image
}