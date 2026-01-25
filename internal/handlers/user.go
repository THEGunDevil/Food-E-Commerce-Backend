package handlers

import (
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetUsersHandler fetches paginated users
func GetUsersHandler(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	params := gen.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// 1️⃣ Fetch paginated users
	users, err := db.Q.ListUsers(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "no user found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "failed to fetch user",
		})
		return
	}

	// 2️⃣ Count total users
	totalCount, err := db.Q.CountUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// 3️⃣ Build response
	response := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		// Convert UUID properly
		userUUID := uuid.UUID(user.ID.Bytes)

		// Handle nullable fields
		var banUntilPtr *time.Time
		if user.BanUntil.Valid {
			banUntilPtr = &user.BanUntil.Time
		}

		var banReasonPtr *string
		if user.BanReason.Valid {
			banReasonPtr = &user.BanReason.String
		}

		response = append(response, models.UserResponse{
			ID:             userUUID,
			Email:          user.Email,
			Phone:          user.Phone,    // Changed from PhoneNumber to Phone
			FullName:       user.FullName, // Combined FirstName + LastName
			Bio:            &user.Bio.String,
			Role:           models.UserRole(user.Role.String),
			AvatarURL:      services.GetStringPtr(user.AvatarUrl), // Changed from ProfileImg
			TokenVersion:   int32(user.TokenVersion),
			IsBanned:       user.IsBanned.Bool,
			BanReason:      banReasonPtr,
			BanUntil:       banUntilPtr,
			IsPermanentBan: user.IsPermanentBan.Bool,
			IsActive:       user.IsActive.Bool,
			IsVerified:     user.IsVerified.Bool,
			LastLogin:      services.GetTimePtr(user.LastLogin),
			CreatedAt:      user.CreatedAt.Time,
			UpdatedAt:      user.UpdatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"limit":       limit,
		"count":       len(response),
		"total_count": totalCount,
		"total_pages": totalPages,
		"users":       response,
	})
}

func GetUserByIDHandler(c *gin.Context) {
	// Get user ID from JWT token (auth middleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Fetch user from DB
	user, err := db.Q.GetUserByID(c.Request.Context(), pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Handle nullable fields
	var banUntilPtr *time.Time
	if user.BanUntil.Valid {
		banUntilPtr = &user.BanUntil.Time
	}

	var banReasonPtr *string
	if user.BanReason.Valid {
		banReasonPtr = &user.BanReason.String
	}

	// Build response
	resp := models.UserResponse{
		ID:             userUUID,
		Email:          user.Email,
		Phone:          user.Phone,
		FullName:       user.FullName,
		Bio:            &user.Bio.String,
		Role:           models.UserRole(user.Role.String),
		AvatarURL:      &user.AvatarUrl.String,
		TokenVersion:   int32(user.TokenVersion),
		IsBanned:       user.IsBanned.Bool,
		BanReason:      banReasonPtr,
		BanUntil:       banUntilPtr,
		IsPermanentBan: user.IsPermanentBan.Bool,
		IsActive:       user.IsActive.Bool,
		IsVerified:     user.IsVerified.Bool,
		LastLogin:      services.GetTimePtr(user.LastLogin),
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}

	log.Printf("👤 Returning user data for user %v (banned: %v)", userUUID, user.IsBanned.Bool)
	c.JSON(http.StatusOK, resp)
}

func SearchUsersPaginatedHandler(c *gin.Context) {
	// Pagination
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

	// Search query
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	log.Printf("🔍 Searching users: query='%s', page=%d", query, page)

	// SQLC params - adjust based on your actual query
	params := gen.SearchUsersParams{
		Column1: pgtype.Text{String: "%" + query + "%", Valid: true}, // Use LIKE pattern
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	// Execute search query
	users, err := db.Q.SearchUsers(c.Request.Context(), params)
	if err != nil {
		log.Printf("❌ Search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}

	// Count total matching users
	totalCount, err := db.Q.CountUsers(c.Request.Context())
	if err != nil {
		log.Printf("❌ Count error: %v", err)
		totalCount = 0
	}

	// Map response
	result := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		userUUID := uuid.UUID(user.ID.Bytes)
		result = append(result, models.UserResponse{
			ID:        userUUID,
			Email:     user.Email,
			Phone:     user.Phone,
			FullName:  user.FullName,
			Role:      models.UserRole(user.Role.String),
			AvatarURL: services.GetStringPtr(user.AvatarUrl),
			IsActive:  user.IsActive.Bool,
			CreatedAt: user.CreatedAt.Time,
		})
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"limit":       limit,
		"count":       len(result),
		"total_count": totalCount,
		"total_pages": totalPages,
		"users":       result,
	})
}

// UpdateUserByIDHandler updates user by ID
func UpdateUserByIDHandler(c *gin.Context) {
	// Parse UUID
	idStr := c.Param("id")
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse request
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start building update params
	params := gen.UpdateUserParams{
		ID: pgtype.UUID{Bytes: parsedID, Valid: true},
	}

	// Update fields if provided
	if req.Email != nil {
		params.Email = pgtype.Text{String: *req.Email, Valid: true}
	}

	if req.Phone != nil {
		params.Phone = pgtype.Text{String: *req.Phone, Valid: true}
	}

	if req.FullName != nil {
		params.FullName = pgtype.Text{String: *req.FullName, Valid: true}
	}

	if req.Bio != nil {
		params.Bio = pgtype.Text{String: *req.Bio, Valid: true}
	}

	if req.Role != nil {
		params.Role = pgtype.Text{String: string(*req.Role), Valid: true}
	}

	if req.IsActive != nil {
		params.IsActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	if req.IsVerified != nil {
		params.IsVerified = pgtype.Bool{Bool: *req.IsVerified, Valid: true}
	}

	if req.AvatarURL != nil {
		// Upload to cloud storage and get URL
	}

	// Save changes
	updatedUser, err := db.Q.UpdateUser(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed: " + err.Error()})
		return
	}

	// Build response
	var banUntilPtr *time.Time
	if updatedUser.BanUntil.Valid {
		banUntilPtr = &updatedUser.BanUntil.Time
	}

	var banReasonPtr *string
	if updatedUser.BanReason.Valid {
		banReasonPtr = &updatedUser.BanReason.String
	}

	resp := models.UserResponse{
		ID:             uuid.UUID(updatedUser.ID.Bytes),
		Email:          updatedUser.Email,
		Phone:          updatedUser.Phone,
		FullName:       updatedUser.FullName,
		Bio:            &updatedUser.Bio.String,
		Role:           models.UserRole(updatedUser.Role.String),
		AvatarURL:      services.GetStringPtr(updatedUser.AvatarUrl),
		TokenVersion:   int32(updatedUser.TokenVersion),
		IsBanned:       updatedUser.IsBanned.Bool,
		BanReason:      banReasonPtr,
		BanUntil:       banUntilPtr,
		IsPermanentBan: updatedUser.IsPermanentBan.Bool,
		IsActive:       updatedUser.IsActive.Bool,
		IsVerified:     updatedUser.IsVerified.Bool,
		LastLogin:      services.GetTimePtr(updatedUser.LastLogin),
		CreatedAt:      updatedUser.CreatedAt.Time,
		UpdatedAt:      updatedUser.UpdatedAt.Time,
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteProfileImage(c *gin.Context) {
	// Get userID from context (from auth middleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Fetch user to check existing avatar
	user, err := db.Q.GetUserByID(c.Request.Context(), services.UUIDToPGType(userUUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// If no avatar exists
	if !user.AvatarUrl.Valid || user.AvatarUrl.String == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No profile image to delete"})
		return
	}

	err = services.DeleteImageFromCloudinary(user.AvatarPublicID.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	params := gen.UpdateUserParams{
		ID:             pgtype.UUID{Bytes: userUUID, Valid: true},
		AvatarUrl:      pgtype.Text{Valid: false}, // Set to NULL
		AvatarPublicID: pgtype.Text{Valid: false},
	}

	_, err = db.Q.UpdateUser(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile image deleted successfully",
	})
}

func BanUserByIDHandler(c *gin.Context) {
	// Parse UUID
	idStr := c.Param("id")
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Bind request
	var req models.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle BanUntil
	var banUntil pgtype.Timestamp
	if req.Permanent {
		// Permanent ban has no expiration
		banUntil = pgtype.Timestamp{Valid: false}
	} else if req.BanUntil != nil {
		banUntil = pgtype.Timestamp{Time: *req.BanUntil, Valid: true}
	} else {
		// Default ban duration (e.g., 30 days)
		defaultBanUntil := time.Now().Add(30 * 24 * time.Hour)
		banUntil = pgtype.Timestamp{Time: defaultBanUntil, Valid: true}
	}

	// Update user ban
	params := gen.BanUserParams{
		ID:             pgtype.UUID{Bytes: parsedID, Valid: true},
		BanReason:      pgtype.Text{String: req.Reason, Valid: req.Reason != ""},
		BanUntil:       banUntil,
		IsPermanentBan: pgtype.Bool{Bool: req.Permanent, Valid: true},
	}

	updatedUser, err := db.Q.BanUser(c.Request.Context(), params)
	if err != nil {
		log.Printf("BanUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user ban status"})
		return
	}

	// Prepare response
	var banUntilPtr *time.Time
	if updatedUser.BanUntil.Valid {
		banUntilPtr = &updatedUser.BanUntil.Time
	}

	var banReasonPtr *string
	if updatedUser.BanReason.Valid {
		banReasonPtr = &updatedUser.BanReason.String
	}

	resp := models.UserResponse{
		ID:             uuid.UUID(updatedUser.ID.Bytes),
		Email:          updatedUser.Email,
		Phone:          updatedUser.Phone,
		FullName:       updatedUser.FullName,
		Role:           models.UserRole(updatedUser.Role.String),
		AvatarURL:      services.GetStringPtr(updatedUser.AvatarUrl),
		TokenVersion:   int32(updatedUser.TokenVersion),
		IsBanned:       updatedUser.IsBanned.Bool,
		BanReason:      banReasonPtr,
		BanUntil:       banUntilPtr,
		IsPermanentBan: updatedUser.IsPermanentBan.Bool,
		IsActive:       updatedUser.IsActive.Bool,
		IsVerified:     updatedUser.IsVerified.Bool,
		CreatedAt:      updatedUser.CreatedAt.Time,
		UpdatedAt:      updatedUser.UpdatedAt.Time,
	}

	c.JSON(http.StatusOK, resp)
}

// Add more handlers as needed:
func UnbanUserByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	updatedUser, err := db.Q.UnbanUser(c.Request.Context(), services.UUIDToPGType(parsedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User unbanned successfully",
		"user_id": updatedUser.ID.Bytes,
	})
}
