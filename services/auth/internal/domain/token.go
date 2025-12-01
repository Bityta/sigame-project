package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// TokenPair represents a pair of access and refresh tokens
// All fields are required
type TokenPair struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterRequest represents the registration request body
// All fields are required
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"player123"`
	Password string `json:"password" binding:"required" example:"securepass123"`
}

// LoginRequest represents the login request body
// All fields are required
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"player123"`
	Password string `json:"password" binding:"required" example:"securepass123"`
}

// RefreshRequest represents the token refresh request body
// RefreshToken is required
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterResponse represents the registration response
// All fields are always present
type RegisterResponse struct {
	User         *UserResponse `json:"user" binding:"required"`
	AccessToken  string        `json:"access_token" binding:"required"`
	RefreshToken string        `json:"refresh_token" binding:"required"`
}

// LoginResponse represents the login response
// All fields are always present
type LoginResponse struct {
	User         *UserResponse `json:"user" binding:"required"`
	AccessToken  string        `json:"access_token" binding:"required"`
	RefreshToken string        `json:"refresh_token" binding:"required"`
}

// CheckUsernameResponse represents the username availability check response
// Available is required, Reason is optional (only present when available=false)
type CheckUsernameResponse struct {
	Available bool   `json:"available" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// ErrorResponse represents an error response
// Error code is required, Message is optional
type ErrorResponse struct {
	Error   string `json:"error" binding:"required"`
	Message string `json:"message,omitempty"`
}

// HashToken creates a SHA-256 hash of a token for storage
// Used for refresh token storage in database
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GetTokenID extracts the token ID (jti) from claims
func GetTokenID(claims *Claims) string {
	return claims.ID
}

