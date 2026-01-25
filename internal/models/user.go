package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type APIResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
	RoleStaff    UserRole = "staff"
)

type User struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	FullName          string     `json:"full_name"`
	Bio               *string    `json:"bio"`
	PasswordHash      string     `json:"-"`
	Role              UserRole   `json:"role"`
	AvatarURL         *string    `json:"avatar_url"`
	AvatarURLPublicID *string    `json:"avatar_public_id"`
	TokenVersion      int32      `json:"token_version"`
	IsBanned          bool       `json:"is_banned"`
	BanReason         *string    `json:"ban_reason,omitempty"`
	BanUntil          *time.Time `json:"ban_until,omitempty"`
	IsPermanentBan    bool       `json:"is_permanent_ban"`
	IsActive          bool       `json:"is_active"`
	IsVerified        bool       `json:"is_verified"`
	LastLogin         *time.Time `json:"last_login"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
type CreateUserRequest struct {
	Email           string  `json:"email" validate:"required,email"`
	Phone           string  `json:"phone" validate:"required"`
	FullName        string  `json:"full_name" validate:"required,min=2,max=100"`
	Bio             *string `json:"bio"`
	Password        string  `json:"password" validate:"required,min=8"`
	ConfirmPassword string  `json:"confirm_password" validate:"required,min=8"`
	AvatarURL       *string `json:"avatar_url"`
}

type UpdateUserRequest struct {
	Email      *string   `json:"email" validate:"omitempty,email"`
	Phone      *string   `json:"phone"`
	FullName   *string   `json:"full_name" validate:"omitempty,min=2,max=100"`
	Bio        *string   `json:"bio"`
	AvatarURL  *string   `json:"avatar_url"`
	Role       *UserRole `json:"role" validate:"omitempty,oneof=customer admin staff"`
	IsActive   *bool     `json:"is_active"`
	IsVerified *bool     `json:"is_verified"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type BanUserRequest struct {
	Reason    string     `json:"reason"`
	BanUntil  *time.Time `json:"ban_until"`
	Permanent bool       `json:"is_permanent"`
}

type UserResponse struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	FullName          string     `json:"full_name"`
	Bio               *string    `json:"bio"`
	Role              UserRole   `json:"role"`
	AvatarURL         *string    `json:"avatar_url"`
	AvatarURLPublicID *string    `json:"avatar_public_id"`
	TokenVersion      int32      `json:"token_version"`
	IsBanned          bool       `json:"is_banned"`
	BanReason         *string    `json:"ban_reason,omitempty"`
	BanUntil          *time.Time `json:"ban_until,omitempty"`
	IsPermanentBan    bool       `json:"is_permanent_ban"`
	IsActive          bool       `json:"is_active"`
	IsVerified        bool       `json:"is_verified"`
	LastLogin         *time.Time `json:"last_login"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (u *User) IsCurrentlyBanned() bool {
	if !u.IsBanned {
		return false
	}

	if u.IsPermanentBan {
		return true
	}

	if u.BanUntil == nil {
		return u.IsBanned
	}

	return u.BanUntil.After(time.Now())
}

func (u *User) CanLogin() bool {
	return u.IsActive && !u.IsCurrentlyBanned()
}
