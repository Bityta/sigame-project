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

// GameRepository handles game state persistence in Redis
type GameRepository struct {
	client *redis.Client
}

// NewGameRepository creates a new GameRepository
func NewGameRepository(client *redis.Client) *GameRepository {
	return &GameRepository{client: client}
}

// SaveGameState saves the current game state to Redis
func (r *GameRepository) SaveGameState(ctx context.Context, game *domain.Game) error {
	key := fmt.Sprintf("game:%s:state", game.ID.String())

	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	ttl := config.GetCacheTTL("game_state")
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	return nil
}

// LoadGameState loads game state from Redis
func (r *GameRepository) LoadGameState(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	key := fmt.Sprintf("game:%s:state", gameID.String())

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

// DeleteGameState deletes game state from Redis
func (r *GameRepository) DeleteGameState(ctx context.Context, gameID uuid.UUID) error {
	key := fmt.Sprintf("game:%s:state", gameID.String())
	return r.client.Del(ctx, key).Err()
}

// SavePlayerScore saves a player's score
func (r *GameRepository) SavePlayerScore(ctx context.Context, gameID, userID uuid.UUID, score int) error {
	key := fmt.Sprintf("game:%s:scores", gameID.String())
	return r.client.HSet(ctx, key, userID.String(), score).Err()
}

// GetPlayerScore retrieves a player's score
func (r *GameRepository) GetPlayerScore(ctx context.Context, gameID, userID uuid.UUID) (int, error) {
	key := fmt.Sprintf("game:%s:scores", gameID.String())
	score, err := r.client.HGet(ctx, key, userID.String()).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return score, err
}

// GetAllScores retrieves all player scores for a game
func (r *GameRepository) GetAllScores(ctx context.Context, gameID uuid.UUID) (map[string]int, error) {
	key := fmt.Sprintf("game:%s:scores", gameID.String())
	scores, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for userID, scoreStr := range scores {
		var score int
		fmt.Sscanf(scoreStr, "%d", &score)
		result[userID] = score
	}

	return result, nil
}

// AddActivePlayer adds a player to the game's active players set
func (r *GameRepository) AddActivePlayer(ctx context.Context, gameID, userID uuid.UUID) error {
	key := fmt.Sprintf("game:%s:players", gameID.String())
	return r.client.SAdd(ctx, key, userID.String()).Err()
}

// RemoveActivePlayer removes a player from the game's active players set
func (r *GameRepository) RemoveActivePlayer(ctx context.Context, gameID, userID uuid.UUID) error {
	key := fmt.Sprintf("game:%s:players", gameID.String())
	return r.client.SRem(ctx, key, userID.String()).Err()
}

// GetActivePlayers retrieves all active player IDs for a game
func (r *GameRepository) GetActivePlayers(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error) {
	key := fmt.Sprintf("game:%s:players", gameID.String())
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	players := make([]uuid.UUID, 0, len(members))
	for _, member := range members {
		if userID, err := uuid.Parse(member); err == nil {
			players = append(players, userID)
		}
	}

	return players, nil
}

// SetGameMetadata saves game metadata
func (r *GameRepository) SetGameMetadata(ctx context.Context, gameID uuid.UUID, metadata map[string]interface{}) error {
	key := fmt.Sprintf("game:%s:meta", gameID.String())
	
	for field, value := range metadata {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case int:
			strValue = fmt.Sprintf("%d", v)
		case bool:
			strValue = fmt.Sprintf("%t", v)
		default:
			data, _ := json.Marshal(v)
			strValue = string(data)
		}
		
		if err := r.client.HSet(ctx, key, field, strValue).Err(); err != nil {
			return err
		}
	}

	// Set TTL
	ttl := config.GetCacheTTL("game_state")
	return r.client.Expire(ctx, key, ttl).Err()
}

// GetGameMetadata retrieves game metadata
func (r *GameRepository) GetGameMetadata(ctx context.Context, gameID uuid.UUID) (map[string]string, error) {
	key := fmt.Sprintf("game:%s:meta", gameID.String())
	return r.client.HGetAll(ctx, key).Result()
}

// SetActiveGame adds game to active games sorted set
func (r *GameRepository) SetActiveGame(ctx context.Context, gameID uuid.UUID, timestamp time.Time) error {
	key := "games:active"
	score := float64(timestamp.Unix())
	return r.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: gameID.String(),
	}).Err()
}

// RemoveActiveGame removes game from active games
func (r *GameRepository) RemoveActiveGame(ctx context.Context, gameID uuid.UUID) error {
	key := "games:active"
	return r.client.ZRem(ctx, key, gameID.String()).Err()
}

// GetActiveGames retrieves all active game IDs
func (r *GameRepository) GetActiveGames(ctx context.Context, limit int64) ([]uuid.UUID, error) {
	key := "games:active"
	members, err := r.client.ZRevRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	games := make([]uuid.UUID, 0, len(members))
	for _, member := range members {
		if gameID, err := uuid.Parse(member); err == nil {
			games = append(games, gameID)
		}
	}

	return games, nil
}

