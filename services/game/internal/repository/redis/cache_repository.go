package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sigame/game/internal/config"
	"github.com/sigame/game/internal/domain"
)

// CacheRepository handles caching operations in Redis
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new CacheRepository
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// CachePack caches a pack in Redis
func (r *CacheRepository) CachePack(ctx context.Context, pack *domain.Pack) error {
	key := fmt.Sprintf("pack:%s:content", pack.ID.String())

	data, err := json.Marshal(pack)
	if err != nil {
		return fmt.Errorf("failed to marshal pack: %w", err)
	}

	ttl := config.GetCacheTTL("pack")
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to cache pack: %w", err)
	}

	return nil
}

// GetCachedPack retrieves a cached pack from Redis
func (r *CacheRepository) GetCachedPack(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	key := fmt.Sprintf("pack:%s:content", packID.String())

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, domain.ErrPackNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached pack: %w", err)
	}

	var pack domain.Pack
	if err := json.Unmarshal(data, &pack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pack: %w", err)
	}

	return &pack, nil
}

// InvalidatePackCache removes a pack from cache
func (r *CacheRepository) InvalidatePackCache(ctx context.Context, packID uuid.UUID) error {
	key := fmt.Sprintf("pack:%s:content", packID.String())
	return r.client.Del(ctx, key).Err()
}

// SetWithTTL sets a key with a custom TTL
func (r *CacheRepository) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value by key
func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Delete deletes a key
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

