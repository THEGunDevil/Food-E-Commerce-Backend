package services

import (
	"fmt"
	"github.com/THEGunDevil/Food-E-Commerce-Backend.git/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype" // FloatToNumeric converts float64 to pgtype.Float8
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UUIDToPGType(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}
func PgtypeToUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	u, _ := uuid.FromBytes(p.Bytes[:])
	return u
}

func StringToPGUUID(s string) pgtype.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// Converts string → pgtype.Text
func StringToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}
func SafeInt(value interface{}) int {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func StringToFloat(s string) float64 {
	if s == "" {
		return 0.0
	}

	// Attempt the conversion
	parseFloat, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return parseFloat
}

func FloatToPGNumeric(f interface{}) pgtype.Numeric {
	switch v := f.(type) {
	case float32:
		return pgtype.Numeric{
			Int:   big.NewInt(int64(v * 100)), // multiply for decimal precision
			Exp:   -2,
			Valid: true,
		}

	case float64:
		return pgtype.Numeric{
			Int:   big.NewInt(int64(v * 100)), // multiply for decimal precision
			Exp:   -2,
			Valid: true,
		}
	default:
		// If type is not supported, return NULL numeric
		return pgtype.Numeric{Valid: false}
	}
}

func GetStringPtr(text pgtype.Text) *string {
	if text.Valid {
		return &text.String
	}
	return nil
}

func GetTimePtr(timestamp pgtype.Timestamp) *time.Time {
	if timestamp.Valid {
		return &timestamp.Time
	}
	return nil
}

func NumericToPtr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, _ := n.Float64Value()
	val := f.Float64
	return &val
}
func NumericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	f, _ := n.Float64Value()
	val := f.Float64
	return val
}
func TimeToTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func PgTypeBool(b *bool) pgtype.Bool {
	if b != nil {
		return pgtype.Bool{Bool: *b, Valid: true}
	}
	return pgtype.Bool{Valid: false}
}

func HandleValidationError(c *gin.Context, err error) {
	var errMsg string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Create a friendly error message list
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, fmt.Sprintf("Field '%s' failed validation: %s", e.Field(), e.Tag()))
		}
		errMsg = strings.Join(errors, ", ")
	} else {
		errMsg = err.Error()
	}

	c.JSON(http.StatusBadRequest, models.APIResponse{
		Success: false,
		Message: "invalid form data",
		Error:   errMsg,
	})
}
