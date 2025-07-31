package middleware

import (
	"time"
)

type AuthSessionDetails struct {
	SessionID int64     `json:"session_id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiredAt time.Time `json:"expired_at"`
}
