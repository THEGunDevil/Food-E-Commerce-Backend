package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description"`
	CatImageUrl      string    `json:"cat_image_url"`
	CatImagePublicID string    `json:"cat_image_public_id"`
	DisplayOrder     int32     `json:"display_order"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateCategoryRequest for creating a new category
type CreateCategoryRequest struct {
	Name         string                `form:"name" validate:"required,min=2,max=100"`
	Slug         string                `form:"slug" validate:"required,min=2,max=100,slug"`
	Description  string                `form:"description" validate:"max=1000"`
	DisplayOrder *int32                `form:"display_order" validate:"omitempty,gte=0"`
	IsActive     *bool                 `form:"is_active"`
	Image        *multipart.FileHeader `form:"image" validate:"required"`
}

// UpdateCategoryRequest for updating an existing category
type UpdateCategoryRequest struct {
	Name             *string               `form:"name" validate:"omitempty,min=2,max=100"`
	Slug             *string               `form:"slug" validate:"omitempty,min=2,max=100,slug"`
	Description      *string               `form:"description" validate:"omitempty,max=1000"`
	DisplayOrder     *int32                `form:"display_order" validate:"omitempty,gte=0"`
	IsActive         *bool                 `form:"is_active"`
	Image            *multipart.FileHeader `form:"image" validate:"omitempty"`
	CatImagePublicID *string               `form:"cat_image_public_id" validate:"omitempty"`
	RemoveImage      *bool                 `form:"remove_image" validate:"omitempty"`
}

// ListCategoriesQuery for filtering/sorting categories
type ListCategoriesQuery struct {
	Page      int    `query:"page" validate:"omitempty,gte=1"`
	PageSize  int    `query:"page_size" validate:"omitempty,gte=1,lte=100"`
	Search    string `query:"search" validate:"omitempty,max=100"`
	IsActive  *bool  `query:"is_active"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=name display_order created_at"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// CategoryListResponse for paginated category list
type CategoryListResponse struct {
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	Categories []Category `json:"categories"`
}
