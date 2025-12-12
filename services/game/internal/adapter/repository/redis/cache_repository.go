package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sigame/game/internal/infrastructure/config"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/domain/event"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) CachePack(ctx context.Context, pack *domain.Pack) error {
	key := packKey(pack.ID)

	data, err := json.Marshal(pack)
	if err != nil {
		return fmt.Errorf("failed to marshal pack: %w", err)
	}

	ttl := config.PackCacheTTL
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to cache pack: %w", err)
	}

	return nil
}

func (r *CacheRepository) GetCachedPack(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	key := packKey(packID)

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

func (r *CacheRepository) InvalidatePackCache(ctx context.Context, packID uuid.UUID) error {
	key := packKey(packID)
	return r.client.Del(ctx, key).Err()
}

func (r *CacheRepository) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

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

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
