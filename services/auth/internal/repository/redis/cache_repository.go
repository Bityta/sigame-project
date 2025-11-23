package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	
	"github.com/sigame/auth/internal/domain"
)

// CacheRepository handles caching operations in Redis
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new Redis cache repository instance
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// CacheUsernameExists caches the result of username existence check
func (r *CacheRepository) CacheUsernameExists(ctx context.Context, username string, exists bool, ttl time.Duration) error {
	key := domain.RedisKeyUsernameExists + username
	value := "0"
	if exists {
		value = "1"
	}
	
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache username existence: %w", err)
	}
	
	return nil
}

// GetUsernameExists retrieves cached username existence check result
// Returns: exists (bool), found (bool), error
func (r *CacheRepository) GetUsernameExists(ctx context.Context, username string) (bool, bool, error) {
	key := domain.RedisKeyUsernameExists + username
	
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, false, nil // Not found in cache
	}
	if err != nil {
		return false, false, fmt.Errorf("failed to get username existence from cache: %w", err)
	}
	
	exists := value == "1"
	return exists, true, nil
}

// AddToBlacklist adds a token to the blacklist
func (r *CacheRepository) AddToBlacklist(ctx context.Context, tokenID string, ttl time.Duration) error {
	key := domain.RedisKeyBlacklist + tokenID
	
	err := r.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	
	return nil
}

// IsBlacklisted checks if a token is blacklisted
func (r *CacheRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := domain.RedisKeyBlacklist + tokenID
	
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	
	return exists > 0, nil
}

// IncrementRateLimit increments the rate limit counter for an IP
// Returns the current count
func (r *CacheRepository) IncrementRateLimit(ctx context.Context, ip string, window time.Duration) (int, error) {
	key := domain.RedisKeyRateLimit + ip
	
	// Use pipeline for atomicity
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment rate limit: %w", err)
	}
	
	count, err := incr.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}
	
	return int(count), nil
}

// GetRateLimitCount gets the current rate limit count for an IP
func (r *CacheRepository) GetRateLimitCount(ctx context.Context, ip string) (int, error) {
	key := domain.RedisKeyRateLimit + ip
	
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}
	
	return count, nil
}

// CacheSession stores user session in Redis
func (r *CacheRepository) CacheSession(ctx context.Context, tokenID string, user *domain.User, ttl time.Duration) error {
	key := domain.RedisKeySession + tokenID
	
	sessionData := map[string]interface{}{
		"user_id":    user.ID.String(),
		"username":   user.Username,
		"expires_at": time.Now().Add(ttl).Unix(),
	}
	
	data, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}
	
	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache session: %w", err)
	}
	
	return nil
}

// GetSession retrieves user session from Redis
func (r *CacheRepository) GetSession(ctx context.Context, tokenID string) (*domain.User, error) {
	key := domain.RedisKeySession + tokenID
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}
	
	userID, err := uuid.Parse(sessionData["user_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}
	
	user := &domain.User{
		ID:       userID,
		Username: sessionData["username"].(string),
	}
	
	return user, nil
}

// DeleteSession removes a session from Redis
func (r *CacheRepository) DeleteSession(ctx context.Context, tokenID string) error {
	key := domain.RedisKeySession + tokenID
	
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	
	return nil
}

// ResetRateLimit resets the rate limit counter for an IP
func (r *CacheRepository) ResetRateLimit(ctx context.Context, ip string) error {
	key := domain.RedisKeyRateLimit + ip
	
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}
	
	return nil
}

// CountActiveSessions returns the number of active sessions in Redis
func (r *CacheRepository) CountActiveSessions(ctx context.Context) (int, error) {
	// Count all keys matching "session:*"
	keys, err := r.client.Keys(ctx, domain.RedisKeySession+"*").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count active sessions: %w", err)
	}
	
	return len(keys), nil
}

