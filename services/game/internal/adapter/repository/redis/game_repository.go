package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sigame/game/internal/config"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/domain/event"
)

type GameRepository struct {
	client *redis.Client
}

func NewGameRepository(client *redis.Client) *GameRepository {
	return &GameRepository{client: client}
}

func (r *GameRepository) SaveGameState(ctx context.Context, game *domain.Game) error {
	key := gameStateKey(game.ID)

	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	ttl := config.GameStateCacheTTL
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	return nil
}

func (r *GameRepository) LoadGameState(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	key := gameStateKey(gameID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, domain.ErrGameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load game state: %w", err)
	}

	var game domain.Game
	if err := json.Unmarshal(data, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	return &game, nil
}

func (r *GameRepository) DeleteGameState(ctx context.Context, gameID uuid.UUID) error {
	key := gameStateKey(gameID)
	return r.client.Del(ctx, key).Err()
}

func (r *GameRepository) SavePlayerScore(ctx context.Context, gameID, userID uuid.UUID, score int) error {
	key := gameScoresKey(gameID)
	return r.client.HSet(ctx, key, userID.String(), score).Err()
}

func (r *GameRepository) GetPlayerScore(ctx context.Context, gameID, userID uuid.UUID) (int, error) {
	key := gameScoresKey(gameID)
	score, err := r.client.HGet(ctx, key, userID.String()).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return score, err
}

func (r *GameRepository) GetAllScores(ctx context.Context, gameID uuid.UUID) (map[string]int, error) {
	key := gameScoresKey(gameID)
	scores, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all scores: %w", err)
	}

	result := make(map[string]int, len(scores))
	for userID, scoreStr := range scores {
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse score for user %s: %w", userID, err)
		}
		result[userID] = score
	}

	return result, nil
}

func (r *GameRepository) AddActivePlayer(ctx context.Context, gameID, userID uuid.UUID) error {
	key := gamePlayersKey(gameID)
	return r.client.SAdd(ctx, key, userID.String()).Err()
}

func (r *GameRepository) RemoveActivePlayer(ctx context.Context, gameID, userID uuid.UUID) error {
	key := gamePlayersKey(gameID)
	return r.client.SRem(ctx, key, userID.String()).Err()
}

func (r *GameRepository) GetActivePlayers(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error) {
	key := gamePlayersKey(gameID)
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active players: %w", err)
	}

	return parseUUIDs(members)
}

func (r *GameRepository) SetGameMetadata(ctx context.Context, gameID uuid.UUID, metadata map[string]interface{}) error {
	key := gameMetadataKey(gameID)

	for field, value := range metadata {
		strValue, err := convertToString(value)
		if err != nil {
			return fmt.Errorf("failed to convert field %s to string: %w", field, err)
		}

		if err := r.client.HSet(ctx, key, field, strValue).Err(); err != nil {
			return fmt.Errorf("failed to set metadata field %s: %w", field, err)
		}
	}

	ttl := config.GameStateCacheTTL
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *GameRepository) GetGameMetadata(ctx context.Context, gameID uuid.UUID) (map[string]string, error) {
	key := gameMetadataKey(gameID)
	return r.client.HGetAll(ctx, key).Result()
}

func (r *GameRepository) SetActiveGame(ctx context.Context, gameID uuid.UUID, timestamp time.Time) error {
	key := activeGamesKey()
	score := float64(timestamp.Unix())
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: gameID.String(),
	}).Err()
}

func (r *GameRepository) RemoveActiveGame(ctx context.Context, gameID uuid.UUID) error {
	key := activeGamesKey()
	return r.client.ZRem(ctx, key, gameID.String()).Err()
}

func (r *GameRepository) GetActiveGames(ctx context.Context, limit int64) ([]uuid.UUID, error) {
	key := activeGamesKey()
	members, err := r.client.ZRevRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active games: %w", err)
	}

	return parseUUIDs(members)
}

func parseUUIDs(strs []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(strs))
	for _, str := range strs {
		id, err := uuid.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID %s: %w", str, err)
		}
		result = append(result, id)
	}
	return result, nil
}

func convertToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to marshal value: %w", err)
		}
		return string(data), nil
	}
}
