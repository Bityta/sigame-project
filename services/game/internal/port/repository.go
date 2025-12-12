package port

import (
	"context"

	"github.com/google/uuid"
	"sigame/game/internal/domain/game"
)

type GameRepository interface {
	CreateGameSession(ctx context.Context, g *game.Game) error
	GetGameSession(ctx context.Context, gameID uuid.UUID) (*game.Game, error)
	UpdateGameSession(ctx context.Context, g *game.Game) error
	GetActiveGameForUser(ctx context.Context, userID uuid.UUID) (*game.Game, error)
}

type GameCache interface {
	SaveGameState(ctx context.Context, g *game.Game) error
	LoadGameState(ctx context.Context, gameID uuid.UUID) (*game.Game, error)
	DeleteGameState(ctx context.Context, gameID uuid.UUID) error
}

