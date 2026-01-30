package middleware

import (
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

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cartID uuid.UUID
		var sessionID uuid.UUID

		sessionIDStr, err := c.Cookie("session_id")

		maxAge := 30 * 24 * 60 * 60 // 30 days in seconds

		if err != nil || sessionIDStr == "" {
			// Create new session
			sessionID = uuid.New()

			row, err := db.Q.CreateCartWithUser(c, gen.CreateCartWithUserParams{
				SessionID: services.UUIDToPGType(sessionID),
				UserID:    pgtype.UUID{Valid: false},
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "failed to create cart session",
					Error:   err.Error(),
				})
				c.Abort()
				return
			}

			cartID, err = uuid.FromBytes(row.ID.Bytes[:])
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "invalid cart id",
				})
				c.Abort()
				return
			}

		} else {
			// Existing session
			sessionID, err = uuid.Parse(sessionIDStr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Message: "invalid session id",
				})
				return
			}

			cartRow, err := db.Q.GetCartBySessionID(c, services.UUIDToPGType(sessionID))
			// if services.PgtypeToUUID(cartRow.SessionID) == sessionID {
			// 	c.JSON(http.StatusConflict,models.APIResponse{
			// 		Success: false,
			// 		Message: "cart already exists",
			// 		Error: err.Error(),
			// 	})
			// 	return 
			// } else
			 if err == pgx.ErrNoRows {
				// Create cart if missing
				createRow, err := db.Q.CreateCartWithUser(c, gen.CreateCartWithUserParams{
					SessionID: services.UUIDToPGType(sessionID),
					UserID:    pgtype.UUID{Valid: false},
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.APIResponse{
						Success: false,
						Message: "failed to create cart",
						Error:   err.Error(),
					})
					c.Abort()
					return
				}

				cartID, err = uuid.FromBytes(createRow.ID.Bytes[:])
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.APIResponse{
						Success: false,
						Message: "invalid cart id",
					})
					c.Abort()
					return
				}
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "failed to get cart",
					Error:   err.Error(),
				})
				c.Abort()
				return
			} else {
				cartID, err = uuid.FromBytes(cartRow.ID.Bytes[:])
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.APIResponse{
						Success: false,
						Message: "invalid cart id",
					})
					c.Abort()
					return
				}
			}
		}

		// Set or refresh the session cookie with sliding expiration
		c.SetCookie(
			"session_id",
			sessionID.String(),
			maxAge,
			"/",
			"",
			false,
			true,
		)

		// ✅ GUARANTEED valid UUID
		c.Set("cart_id", cartID)
		c.Next()
	}
}
