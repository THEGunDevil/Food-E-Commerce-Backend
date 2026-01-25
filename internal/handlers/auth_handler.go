package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db"
	gen "github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/db/gen"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/services"
)

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	if len(req.FullName) < 3 || len(req.FullName) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "first and last names must be 3-25 chars"})
		return
	}

	emailRegex := regexp.MustCompile(`^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	if len(req.Email) == 0 || len(req.Email) > 255 || !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	hashed, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user, err := db.Q.CreateUser(c.Request.Context(), gen.CreateUserParams{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Bio:          services.StringToPGText(*req.AvatarURL),
		PasswordHash: hashed,
		AvatarUrl: services.StringToPGText(*req.AvatarURL),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	resp := services.ToUserResponse(user)

	c.JSON(http.StatusCreated, resp)
}

// LoginHandler handles user login and sets refresh token cookie
func LoginHandler(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := db.Q.GetUserByEmail(c, body.Email)
	if err != nil || services.CheckPassword(body.Password, user.PasswordHash) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := services.GenerateAccessToken(user.ID.String(), user.Role.String, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := services.GenerateRefreshToken(user.ID.String(), user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Determine if running on localhost based on Host header
	isLocalhost := strings.Contains(c.Request.Host, "localhost") || strings.Contains(c.Request.Host, "127.0.0.1")
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",                  // Ensure cookie is available site-wide
		MaxAge:   3600 * 24 * 7,        // 7 days
		HttpOnly: true,                 // Prevent JavaScript access
		Secure:   !isLocalhost,         // Secure=true in production, false on localhost
		SameSite: http.SameSiteLaxMode, // Default to Lax for compatibility
	}

	// Use SameSite=None for cross-origin requests in production
	if !isLocalhost {
		cookie.SameSite = http.SameSiteNoneMode
	}

	// Avoid setting Domain explicitly unless necessary
	// If backend is on a different domain (e.g., api.himel-s-library.vercel.app), uncomment and set:
	// cookie.Domain = "your-backend-domain.com"

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"role":         user.Role.String,
	})
}

// RefreshHandler refreshes access token and renews refresh token cookie
func RefreshHandler(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	token, err := services.VerifyToken(cookie, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userIDStr, ok1 := claims["sub"].(string)
	version, ok2 := claims["token_version"].(float64)
	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := db.Q.GetUserByID(c, pgtype.UUID{Bytes: userUUID, Valid: true})
	if err != nil || user.TokenVersion != int32(version) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or invalid"})
		return
	}

	accessToken, err := services.GenerateAccessToken(userIDStr, user.Role.String, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	// Renew refresh token
	refreshToken, err := services.GenerateRefreshToken(userIDStr, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Determine if running on localhost
	isLocalhost := strings.Contains(c.Request.Host, "localhost") || strings.Contains(c.Request.Host, "127.0.0.1")
	cookieConfig := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
		Secure:   !isLocalhost,
		SameSite: http.SameSiteLaxMode,
	}

	if !isLocalhost {
		cookieConfig.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, cookieConfig)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// LogoutHandler clears the refresh token cookie and invalidates token version
func LogoutHandler(c *gin.Context) {
	// Determine if running on localhost
	isLocalhost := strings.Contains(c.Request.Host, "localhost") || strings.Contains(c.Request.Host, "127.0.0.1")
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HttpOnly: true,
		Secure:   !isLocalhost,
		SameSite: http.SameSiteLaxMode,
	}

	if !isLocalhost {
		cookie.SameSite = http.SameSiteNoneMode
	}

	// Invalidate token version for security
	cookieValue, err := c.Cookie("refresh_token")
	if err == nil {
		token, err := services.VerifyToken(cookieValue, true)
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userIDStr, _ := claims["sub"].(string)
				userUUID, _ := uuid.Parse(userIDStr)
				if err := db.Q.UpdateTokenVersion(c.Request.Context(), services.UUIDToPGType(userUUID)); err != nil {
					// Log error but don't fail logout
					log.Printf("Failed to increment token version: %v", err)
				}
			}
		}
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
