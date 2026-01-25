package services

import (
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// BuildCreateMenuItemParams prepares the params for the main menu item insert
func BuildCreateMenuItemParams(
    categoryID string, 
    req *models.CreateMenuItemMultipartRequest, 
    uploadedImages []models.MenuItemImage, // <--- Added this argument back
) (gen.CreateMenuItemParams, []gen.CreateMenuItemImageParams, error) { // <--- Updated return type

    catUUID, err := uuid.Parse(categoryID)
    if err != nil {
        return gen.CreateMenuItemParams{}, nil, err
    }

    params := gen.CreateMenuItemParams{
        CategoryID:   UUIDToPGType(catUUID),
        Name:         req.Name,
        Slug:         req.Slug,
        Description:  StringToPGText(req.Description),
        Price:        FloatToPGNumeric(StringToFloat(req.Price)),
        Ingredients:  req.Ingredients,
        Tags:         req.Tags,
        IsVegetarian: PgTypeBool(req.IsVegetarian),
        IsSpecial:    PgTypeBool(req.IsSpecial),
        IsAvailable:  PgTypeBool(req.IsAvailable),
    }

    if req.DiscountPrice != nil {
        params.DiscountPrice = FloatToPGNumeric(StringToFloat(*req.DiscountPrice))
    }
    if req.PrepTime != nil {
        params.PrepTime = pgtype.Int4{Int32: int32(*req.PrepTime), Valid: true}
    }
    if req.SpicyLevel != nil {
        params.SpicyLevel = pgtype.Int4{Int32: int32(*req.SpicyLevel), Valid: true}
    }
    if req.StockQuantity != nil {
        params.StockQuantity = pgtype.Int4{Int32: int32(*req.StockQuantity), Valid: true}
    }
    if req.MinStockAlert != nil {
        params.MinStockAlert = pgtype.Int4{Int32: int32(*req.MinStockAlert), Valid: true}
    }
    if req.DisplayOrder != nil {
        params.DisplayOrder = pgtype.Int4{Int32: int32(*req.DisplayOrder), Valid: true}
    }

    // Process Images
    imageParams := make([]gen.CreateMenuItemImageParams, len(uploadedImages))
    for i, img := range uploadedImages {
        imageParams[i] = gen.CreateMenuItemImageParams{
            MenuItemID:    UUIDToPGType(img.ID), 
            ImageUrl:      img.ImageUrl,
            ImagePublicID: StringToPGText(img.ImagePublicID),
            IsPrimary:     pgtype.Bool{Bool: i == 0, Valid: true},
            DisplayOrder:  pgtype.Int4{Int32: int32(img.DisplayOrder), Valid: true},
        }
    }

    return params, imageParams, nil
}
// BuildUpdateMenuItemParams constructs UPDATE database parameters from request
func BuildUpdateMenuItemParams(categoryID string, req *models.UpdateMenuItemMultipartRequest) (gen.UpdateMenuItemParams, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return gen.UpdateMenuItemParams{}, err
	}

	params := gen.UpdateMenuItemParams{
		CategoryID: UUIDToPGType(catUUID),
	}

	if req.Name != nil {
		params.Name = *req.Name
	}
	if req.Slug != nil {
		params.Slug = *req.Slug
	}
	if req.Description != nil {
		params.Description = StringToPGText(*req.Description)
	}
	if req.Price != nil {
		params.Price = FloatToPGNumeric(StringToFloat(*req.Price))
	}
	if req.DiscountPrice != nil {
		params.DiscountPrice = FloatToPGNumeric(StringToFloat(*req.DiscountPrice))
	}

	// Always update arrays (if provided in multipart they might be overwrite or append,
	// assuming overwrite based on standard HTTP PUT/PATCH behavior)
	params.Ingredients = req.Ingredients
	params.Tags = req.Tags

	if req.PrepTime != nil {
		params.PrepTime = pgtype.Int4{Int32: int32(*req.PrepTime), Valid: true}
	}
	if req.SpicyLevel != nil {
		params.SpicyLevel = pgtype.Int4{Int32: int32(*req.SpicyLevel), Valid: true}
	}

	params.IsVegetarian = PgTypeBool(req.IsVegetarian)
	params.IsSpecial = PgTypeBool(req.IsSpecial)
	params.IsAvailable = PgTypeBool(req.IsAvailable)

	if req.StockQuantity != nil {
		params.StockQuantity = pgtype.Int4{Int32: int32(*req.StockQuantity), Valid: true}
	}
	if req.MinStockAlert != nil {
		params.MinStockAlert = pgtype.Int4{Int32: int32(*req.MinStockAlert), Valid: true}
	}
	if req.DisplayOrder != nil {
		params.DisplayOrder = pgtype.Int4{Int32: int32(*req.DisplayOrder), Valid: true}
	}

	return params, nil
}
