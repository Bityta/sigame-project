package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	domainGame "sigame/game/internal/domain/game"
	"sigame/game/internal/domain/player"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) CreateGameSession(ctx context.Context, game *domainGame.Game) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, queryInsertGameSession,
		game.ID,
		game.RoomID,
		game.PackID,
		game.Status,
		game.CurrentRound,
		game.CurrentPhase,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create game session: %w", err)
	}

	for _, player := range game.Players {
		if err := r.createGamePlayer(ctx, game.ID, player); err != nil {
			return err
		}
	}

	return nil
}

func (r *GameRepository) UpdateGameSession(ctx context.Context, game *domainGame.Game) error {
	_, err := r.db.ExecContext(ctx, queryUpdateGameSession,
		game.Status,
		game.CurrentRound,
		game.CurrentPhase,
		game.StartedAt,
		game.FinishedAt,
		time.Now(),
		game.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update game session: %w", err)
	}

	for _, player := range game.Players {
		if err := r.updatePlayerScore(ctx, game.ID, player); err != nil {
			return err
		}
	}

	return nil
}

func (r *GameRepository) GetGameSession(ctx context.Context, gameID uuid.UUID) (*domainGame.Game, error) {
	g := &domainGame.Game{
		Players: make(map[uuid.UUID]*player.Player),
	}

	row := r.db.QueryRowContext(ctx, querySelectGameSession, gameID)
	if err := scanGameRow(row, g); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get game session: %w", err)
	}

	if err := loadGamePlayers(ctx, r.db, gameID, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (r *GameRepository) SaveFinalResults(ctx context.Context, game *domainGame.Game) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	_, err = tx.ExecContext(ctx, queryUpdateGameSessionFinal, domainGame.StatusFinished, now, now, game.ID)
	if err != nil {
		return fmt.Errorf("failed to update game session: %w", err)
	}

	for _, player := range game.Players {
		var leftAt *time.Time
		if !player.IsActive {
			leftAt = player.LeftAt
		}

		_, err = tx.ExecContext(ctx, queryUpdatePlayerFinal, player.Score, player.IsActive, leftAt, game.ID, player.UserID)
		if err != nil {
			return fmt.Errorf("failed to update player score: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *GameRepository) createGamePlayer(ctx context.Context, gameID uuid.UUID, player *player.Player) error {
	_, err := r.db.ExecContext(ctx, queryInsertGamePlayer,
		gameID,
		player.UserID,
		player.Username,
		player.Role,
		player.Score,
		player.IsActive,
		player.JoinedAt,
	)

	return err
}

func (r *GameRepository) updatePlayerScore(ctx context.Context, gameID uuid.UUID, player *player.Player) error {
	_, err := r.db.ExecContext(ctx, queryUpdatePlayerScore, player.Score, player.IsActive, gameID, player.UserID)
	return err
}

func (r *GameRepository) GetGamesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domainGame.Game, error) {
	rows, err := r.db.QueryContext(ctx, querySelectGamesByRoomID, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer rows.Close()

	var games []*domainGame.Game
	for rows.Next() {
		g := &domainGame.Game{
			Players: make(map[uuid.UUID]*player.Player),
		}

		if err := scanGame(rows, g); err != nil {
			return nil, err
		}

		games = append(games, g)
	}

	return games, nil
}

func (r *GameRepository) GetActiveGameForUser(ctx context.Context, userID uuid.UUID) (*domainGame.Game, error) {
	g := &domainGame.Game{
		Players: make(map[uuid.UUID]*player.Player),
	}

	row := r.db.QueryRowContext(ctx, querySelectActiveGameForUser, userID)
	if err := scanGameRow(row, g); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get active game for user: %w", err)
	}

	if err := loadGamePlayers(ctx, r.db, g.ID, g); err != nil {
		return nil, err
	}

	return g, nil
}

