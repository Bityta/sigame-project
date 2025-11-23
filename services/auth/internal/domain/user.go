package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the authentication system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RefreshToken represents a refresh token stored in the database
// Refresh tokens are long-lived tokens used to obtain new access tokens
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"` // Hashed token for security
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// UserResponse represents user data safe for API responses
// This excludes sensitive information like password hashes
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to a UserResponse (safe for external use)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

// NormalizeUsername converts username to lowercase for consistent storage
// This prevents duplicate accounts with different casing (e.g., "User" vs "user")
func NormalizeUsername(username string) string {
	return strings.ToLower(username)
}

