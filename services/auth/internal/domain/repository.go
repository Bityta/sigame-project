package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *User) error
	
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	
	// UsernameExists checks if a username already exists
	UsernameExists(ctx context.Context, username string) (bool, error)
	
	// CreateRefreshToken stores a refresh token in the database
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	
	// GetRefreshToken retrieves a refresh token by its hash
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	
	// DeleteRefreshToken removes a refresh token from the database
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	
	// DeleteExpiredTokens removes expired refresh tokens (cleanup job)
	DeleteExpiredTokens(ctx context.Context) (int64, error)
	
	// DeleteUserRefreshTokens deletes all refresh tokens for a user (used during logout)
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
	
	// CountUsers returns the total number of registered users
	CountUsers(ctx context.Context) (int, error)
}

// CacheRepository defines the interface for cache operations
type CacheRepository interface {
	// CacheUsernameExists caches the result of username existence check
	CacheUsernameExists(ctx context.Context, username string, exists bool, ttl time.Duration) error
	
	// GetUsernameExists retrieves cached username existence check result
	// Returns: exists (bool), found (bool), error
	GetUsernameExists(ctx context.Context, username string) (bool, bool, error)
	
	// AddToBlacklist adds a token to the blacklist
	AddToBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error
	
	// IsBlacklisted checks if a token is blacklisted
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
	
	// IncrementRateLimit increments the rate limit counter for an IP
	// Returns the current count
	IncrementRateLimit(ctx context.Context, ip string, window time.Duration) (int, error)
	
	// GetRateLimitCount gets the current rate limit count for an IP
	GetRateLimitCount(ctx context.Context, ip string) (int, error)
	
	// ResetRateLimit resets the rate limit counter for an IP
	ResetRateLimit(ctx context.Context, ip string) error
	
	// CacheSession stores user session in Redis
	CacheSession(ctx context.Context, tokenID string, user *User, ttl time.Duration) error
	
	// GetSession retrieves user session from Redis
	GetSession(ctx context.Context, tokenID string) (*User, error)
	
	// DeleteSession removes a session from Redis
	DeleteSession(ctx context.Context, tokenID string) error
	
	// CountActiveSessions returns the number of active sessions in Redis
	CountActiveSessions(ctx context.Context) (int, error)
}

