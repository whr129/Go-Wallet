package dto

import (
	"time"
)

type CreateUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role"`
}

type UserResponse struct {
	ID                int64     `json:"id"`
	UserName          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             int64        `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

type AuthSessionDetails struct {
	SessionID int64     `json:"session_id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiredAt time.Time `json:"expired_at"`
}
