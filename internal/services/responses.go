package services

import (
	"time"

	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/google/uuid"
)

func ToUserResponse(user gen.User) models.UserResponse {
	return models.UserResponse{
		ID:             user.ID.Bytes,
		Email:          user.Email,
		Phone:          user.Phone,
		FullName:       user.FullName,
		Bio:            &user.Bio.String,
		Role:           models.UserRole(user.Role.String),
		AvatarURL:      &user.AvatarUrl.String,
		TokenVersion:   user.TokenVersion,
		IsBanned:       user.IsBanned.Bool,
		BanReason:      &user.BanReason.String,
		BanUntil:       &user.BanUntil.Time,
		IsPermanentBan: user.IsPermanentBan.Bool,
		IsActive:       user.IsActive.Bool,
		IsVerified:     user.IsVerified.Bool,
		LastLogin:      &user.LastLogin.Time,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}
}
func ToUserAddressResponse(user gen.UserAddress) models.UserAddress {
	var lat float64
	if user.Latitude.Scan(&lat) != nil {
		lat = 0.0 // Set a default value or handle the error appropriately
	}

	var lon float64
	if user.Longitude.Scan(&lon) != nil {
		lon = 0.0 // Set a default value or handle the error appropriately
	}

	// Assign the unmarshalled float64 values to the models.UserAddress fields
	return models.UserAddress{
		ID:           user.ID.Bytes,
		UserID:       user.UserID.Bytes,
		Label:        user.Label,
		AddressLine1: user.AddressLine1,
		AddressLine2: &user.AddressLine2.String,
		Area:         user.Area,
		City:         user.City,
		PostalCode:   &user.PostalCode.String,
		Latitude:     &lat, // Use the correctly scanned 'lat'
		Longitude:    &lon, // Use the correctly scanned 'lon'
		IsDefault:    user.IsDefault.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}
}

func ToMenuItemImage(dbImg gen.MenuItemImage) models.MenuItemImage {
	return models.MenuItemImage{
		ID:            dbImg.ID.Bytes,
		ImageUrl:      dbImg.ImageUrl,
		ImagePublicID: dbImg.ImagePublicID.String,
		IsPrimary:     dbImg.IsPrimary.Bool,
		DisplayOrder:  int(dbImg.DisplayOrder.Int32),
	}
}
func ToMenuItemImages(dbImages []gen.MenuItemImage) []models.MenuItemImage {
	images := make([]models.MenuItemImage, len(dbImages))
	for i, img := range dbImages {
		images[i] = ToMenuItemImage(img)
	}
	return images
}

func ToMenuResponse(
	m gen.MenuItem,
	images []gen.MenuItemImage,
) models.MenuItem {
	return models.MenuItem{
		ID:            m.ID.Bytes,
		CategoryID:    m.CategoryID.Bytes,
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description.String,
		Price:         NumericToFloat(m.Price),
		DiscountPrice: NumericToPtr(m.DiscountPrice),
		Images:        ToMenuItemImages(images),
		Ingredients:   m.Ingredients,
		Tags:          m.Tags,
		PrepTime:      int(m.PrepTime.Int32),
		SpicyLevel:    int(m.SpicyLevel.Int32),
		IsVegetarian:  m.IsVegetarian.Bool,
		IsSpecial:     m.IsSpecial.Bool,
		IsAvailable:   m.IsAvailable.Bool,
		StockQuantity: int(m.StockQuantity.Int32),
		MinStockAlert: int(m.MinStockAlert.Int32),
		TotalOrders:   int(m.TotalOrders.Int32),
		AverageRating: NumericToFloat(m.AverageRating),
		DisplayOrder:  int(m.DisplayOrder.Int32),
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}
func ToMenuListResponseWithCategoryName(
	m gen.ListMenuItemsRow,
	images []gen.MenuItemImage,
) models.MenuItem {
	return models.MenuItem{
		ID:            m.ID.Bytes,
		CategoryName:  m.Categoryname,
		CategoryID:    m.CategoryID.Bytes,
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description.String,
		Price:         NumericToFloat(m.Price),
		DiscountPrice: NumericToPtr(m.DiscountPrice),
		Images:        ToMenuItemImages(images),
		Ingredients:   m.Ingredients,
		Tags:          m.Tags,
		PrepTime:      int(m.PrepTime.Int32),
		SpicyLevel:    int(m.SpicyLevel.Int32),
		IsVegetarian:  m.IsVegetarian.Bool,
		IsSpecial:     m.IsSpecial.Bool,
		IsAvailable:   m.IsAvailable.Bool,
		StockQuantity: int(m.StockQuantity.Int32),
		MinStockAlert: int(m.MinStockAlert.Int32),
		TotalOrders:   int(m.TotalOrders.Int32),
		AverageRating: NumericToFloat(m.AverageRating),
		DisplayOrder:  int(m.DisplayOrder.Int32),
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
	}
}
func ToCategoryResponse(c gen.Category) models.Category {
	return models.Category{
		ID:               c.ID.Bytes,
		Name:             c.Name,
		Slug:             c.Slug,
		Description:      c.Description.String,
		CatImageUrl:      c.CatImageUrl.String,
		CatImagePublicID: c.CatImagePublicID.String,
		DisplayOrder:     c.DisplayOrder.Int32,
		IsActive:         c.IsActive.Bool,
		CreatedAt:        c.CreatedAt.Time,
		UpdatedAt:        c.UpdatedAt.Time,
	}
}
func ToCartItemResponse(c gen.GetMenuItemByIDRow, i gen.CartItem, mis []gen.MenuItemImage) models.CartItemResponse {
	return models.CartItemResponse{
		CartItemID:          i.CartID.Bytes,
		MenuItemID:          PgtypeToUUID(c.ID),
		Name:                c.Name,
		Price:               NumericToFloat(c.Price),
		OriginalPrice:       NumericToFloat(c.DiscountPrice),
		Quantity:            int(i.Quantity),
		Image:               ToMenuItemImages(mis),
		SpecialInstructions: i.SpecialInstructions.String,
	}
}
func ToDeliveryZoneResponse(d gen.DeliveryZone) models.DeliveryZone {
	return models.DeliveryZone{
		ID:              d.ID.Bytes,
		ZoneName:        d.ZoneName,
		AreaNames:       d.AreaNames,
		DeliveryFee:     float64(d.DeliveryFee.Exp),
		MinDeliveryTime: int(d.MinDeliveryTime.Int32),
		MaxDeliveryTime: int(d.MaxDeliveryTime.Int32),
		IsActive:        d.IsActive.Bool,
		CreatedAt:       d.CreatedAt.Time,
		UpdatedAt:       d.UpdatedAt.Time,
	}
}
func ToPromotionResponse(p gen.Promotion) models.Promotion {
	return models.Promotion{
		ID:             p.ID.Bytes,
		Title:          p.Title,
		Description:    p.Description.String,
		DiscountType:   p.DiscountType.String,
		DiscountValue:  float64(p.DiscountValue.Exp),
		MinOrderAmount: float64(p.MinOrderAmount.Exp),
		ValidFrom:      p.ValidFrom.Time,
		ValidUntil:     p.ValidUntil.Time,
		MaxUses:        p.MaxUses.Int32,
		UsedCount:      p.UsedCount.Int32,
		IsActive:       p.IsActive.Bool,
		CreatedAt:      p.CreatedAt.Time,
		UpdatedAt:      p.UpdatedAt.Time,
	}
}

func ToReviewResponse(r gen.Review) models.Review {
	return models.Review{
		ID:     r.ID.Bytes,
		UserID: r.UserID.Bytes,
		// OrderID:    r.OrderID.Bytes,
		MenuItemID: r.MenuItemID.Bytes,
		Rating:     r.Rating,
		Comment:    r.Comment.String,
		// Images:     r.Images,
		IsApproved: r.IsApproved.Bool,
		CreatedAt:  r.CreatedAt.Time,
		UpdatedAt:  r.UpdatedAt.Time,
	}
}
func ToOrderItemResponse(o gen.OrderItem) models.OrderItem {
	return models.OrderItem{
		ID:                  o.ID.Bytes,
		OrderID:             o.OrderID.Bytes,
		MenuItemID:          o.MenuItemID.Bytes,
		MenuItemName:        o.MenuItemName,
		Quantity:            o.Quantity,
		UnitPrice:           float64(o.UnitPrice.Exp),
		TotalPrice:          float64(o.TotalPrice.Exp),
		SpecialInstructions: &o.SpecialInstructions.String,
		CreatedAt:           o.CreatedAt.Time,
		UpdatedAt:           o.UpdatedAt.Time,
	}
}
func ToOrderResponse(o gen.Order) models.Order {
	var userID *uuid.UUID
	if o.UserID.Valid {
		u, _ := uuid.FromBytes(o.UserID.Bytes[:])
		userID = &u
	}

	var deliveryAddressID *uuid.UUID
	if o.DeliveryAddressID.Valid {
		u, _ := uuid.FromBytes(o.DeliveryAddressID.Bytes[:])
		deliveryAddressID = &u
	}

	var deliveryPersonID *uuid.UUID
	if o.DeliveryPersonID.Valid {
		u, _ := uuid.FromBytes(o.DeliveryPersonID.Bytes[:])
		deliveryPersonID = &u
	}

	var customerEmail *string
	if o.CustomerEmail.Valid {
		customerEmail = &o.CustomerEmail.String
	}

	var transactionID *string
	if o.TransactionID.Valid {
		transactionID = &o.TransactionID.String
	}

	var specialInstructions *string
	if o.SpecialInstructions.Valid {
		specialInstructions = &o.SpecialInstructions.String
	}

	var cancelledReason *string
	if o.CancelledReason.Valid {
		cancelledReason = &o.CancelledReason.String
	}

	var estimatedDelivery *time.Time
	if o.EstimatedDelivery.Valid {
		estimatedDelivery = &o.EstimatedDelivery.Time
	}

	var actualDelivery *time.Time
	if o.ActualDelivery.Valid {
		actualDelivery = &o.ActualDelivery.Time
	}

	return models.Order{
		ID:                  PgtypeToUUID(o.ID),
		OrderNumber:         o.OrderNumber,
		UserID:              userID,
		DeliveryAddressID:   deliveryAddressID,
		DeliveryType:        o.DeliveryType.String,
		DeliveryAddress:     o.DeliveryAddress.String,
		CustomerName:        o.CustomerName,
		CustomerPhone:       o.CustomerPhone,
		CustomerEmail:       customerEmail,
		Subtotal:            float64(o.Subtotal.Exp),
		DiscountAmount:      float64(o.DiscountAmount.Exp),
		DeliveryFee:         float64(o.DeliveryFee.Exp),
		VATAmount:           float64(o.VatAmount.Exp),
		TotalAmount:         float64(o.TotalAmount.Exp),
		PaymentMethod:       o.PaymentMethod.String,
		PaymentStatus:       o.PaymentStatus.String,
		TransactionID:       transactionID,
		OrderStatus:         o.OrderStatus.String,
		DeliveryPersonID:    deliveryPersonID,
		EstimatedDelivery:   estimatedDelivery,
		ActualDelivery:      actualDelivery,
		SpecialInstructions: specialInstructions,
		CancelledReason:     cancelledReason,
		CreatedAt:           o.CreatedAt.Time,
		UpdatedAt:           o.UpdatedAt.Time,
	}
}

func ToNotificationResponse(n gen.Notification) models.Notification {
	return models.Notification{
		ID:        n.ID.Bytes,
		UserID:    n.UserID.Bytes,
		EventID:   n.EventID.Bytes,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type.String,
		Priority:  n.Priority.String,
		IsRead:    n.IsRead.Bool,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt.Time,
		UpdatedAt: n.UpdatedAt.Time,
	}
}

// services/event_mapper.go
func ToEventResponse(e gen.Event) models.Event {
	return models.Event{
		ID:          e.ID.Bytes,
		EventType:   e.EventType,
		Payload:     e.Payload,
		Delivered:   e.Delivered.Bool,
		CreatedAt:   e.CreatedAt.Time,
		DeliveredAt: &e.DeliveredAt.Time,
	}
}

func ToFavResponse(f gen.Favorite) models.FavoriteResponse {
	return models.FavoriteResponse{
		ID:         f.ID.Bytes,
		UserID:     f.UserID.Bytes,
		MenuItemID: f.MenuItemID.Bytes,
		CreatedAt:  f.CreatedAt.Time,
	}
}
func ToFavListResponse(f gen.ListFavoritesByUserRow) models.FavoriteResponse {
	return models.FavoriteResponse{
		ID:         f.ID.Bytes,
		MenuItemID: f.MenuItemID.Bytes,
		CreatedAt:  f.CreatedAt.Time,
	}
}
