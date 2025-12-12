package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"sigame/game/internal/infrastructure/config"
	"sigame/game/internal/domain/pack"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) CachePack(ctx context.Context, pack *pack.Pack) error {
	key := packKey(pack.ID)

	data, err := json.Marshal(pack)
	if err != nil {
		return ErrMarshalPack(err)
	}

	ttl := config.PackCacheTTL
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return ErrCachePack(err)
	}

	return nil
}

func (r *CacheRepository) GetCachedPack(ctx context.Context, packID uuid.UUID) (*pack.Pack, error) {
	key := packKey(packID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrPackNotFound
	}
	if err != nil {
		return nil, ErrGetCachedPack(err)
	}

	var pack pack.Pack
	if err := json.Unmarshal(data, &pack); err != nil {
		return nil, ErrUnmarshalPack(err)
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
		return ErrMarshalValue(err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return ErrKeyNotFound(key)
	}
	if err != nil {
		return ErrGetValue(err)
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
