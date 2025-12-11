package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
)

// GameRepository handles game persistence in PostgreSQL
type GameRepository struct {
	db *sql.DB
}

// NewGameRepository creates a new GameRepository
func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

// CreateGameSession creates a new game session record
func (r *GameRepository) CreateGameSession(ctx context.Context, game *domain.Game) error {
	query := `
		INSERT INTO game_sessions (id, room_id, pack_id, status, current_round, current_phase, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
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

	// Insert players
	for _, player := range game.Players {
		if err := r.createGamePlayer(ctx, game.ID, player); err != nil {
			return err
		}
	}

	return nil
}

// UpdateGameSession updates an existing game session
func (r *GameRepository) UpdateGameSession(ctx context.Context, game *domain.Game) error {
	query := `
		UPDATE game_sessions 
		SET status = $1, current_round = $2, current_phase = $3, 
		    started_at = $4, finished_at = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
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

	// Update player scores
	for _, player := range game.Players {
		if err := r.updatePlayerScore(ctx, game.ID, player); err != nil {
			return err
		}
	}

	return nil
}

// GetGameSession retrieves a game session by ID
func (r *GameRepository) GetGameSession(ctx context.Context, gameID uuid.UUID) (*domain.Game, error) {
	query := `
		SELECT id, room_id, pack_id, status, current_round, current_phase, 
		       started_at, finished_at, created_at, updated_at
		FROM game_sessions
		WHERE id = $1
	`

	game := &domain.Game{
		Players: make(map[uuid.UUID]*domain.Player),
	}

	var startedAt, finishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, gameID).Scan(
		&game.ID,
		&game.RoomID,
		&game.PackID,
		&game.Status,
		&game.CurrentRound,
		&game.CurrentPhase,
		&startedAt,
		&finishedAt,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrGameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get game session: %w", err)
	}

	if startedAt.Valid {
		game.StartedAt = &startedAt.Time
	}
	if finishedAt.Valid {
		game.FinishedAt = &finishedAt.Time
	}

	// Load players
	players, err := r.getGamePlayers(ctx, gameID)
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		game.Players[player.UserID] = player
	}

	return game, nil
}

// SaveFinalResults saves final game results
func (r *GameRepository) SaveFinalResults(ctx context.Context, game *domain.Game) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update game session
	query := `
		UPDATE game_sessions 
		SET status = $1, finished_at = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now()
	_, err = tx.ExecContext(ctx, query, domain.GameStatusFinished, now, now, game.ID)
	if err != nil {
		return fmt.Errorf("failed to update game session: %w", err)
	}

	// Update all player scores
	for _, player := range game.Players {
		query = `
			UPDATE game_players 
			SET score = $1, is_active = $2, left_at = $3
			WHERE game_id = $4 AND user_id = $5
		`

		var leftAt *time.Time
		if !player.IsActive {
			leftAt = player.LeftAt
		}

		_, err = tx.ExecContext(ctx, query, player.Score, player.IsActive, leftAt, game.ID, player.UserID)
		if err != nil {
			return fmt.Errorf("failed to update player score: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// createGamePlayer inserts a player record
func (r *GameRepository) createGamePlayer(ctx context.Context, gameID uuid.UUID, player *domain.Player) error {
	query := `
		INSERT INTO game_players (game_id, user_id, username, role, score, is_active, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
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

// updatePlayerScore updates a player's score
func (r *GameRepository) updatePlayerScore(ctx context.Context, gameID uuid.UUID, player *domain.Player) error {
	query := `
		UPDATE game_players 
		SET score = $1, is_active = $2
		WHERE game_id = $3 AND user_id = $4
	`

	_, err := r.db.ExecContext(ctx, query, player.Score, player.IsActive, gameID, player.UserID)
	return err
}

// getGamePlayers retrieves all players for a game
func (r *GameRepository) getGamePlayers(ctx context.Context, gameID uuid.UUID) ([]*domain.Player, error) {
	query := `
		SELECT user_id, username, role, score, is_active, joined_at, left_at
		FROM game_players
		WHERE game_id = $1
		ORDER BY joined_at
	`

	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game players: %w", err)
	}
	defer rows.Close()

	var players []*domain.Player
	for rows.Next() {
		player := &domain.Player{}
		var role string
		var leftAt sql.NullTime

		err := rows.Scan(
			&player.UserID,
			&player.Username,
			&role,
			&player.Score,
			&player.IsActive,
			&player.JoinedAt,
			&leftAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}

		player.Role = domain.PlayerRole(role)
		if leftAt.Valid {
			player.LeftAt = &leftAt.Time
		}

		players = append(players, player)
	}

	return players, nil
}

// GetGamesByRoomID retrieves games by room ID
func (r *GameRepository) GetGamesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Game, error) {
	query := `
		SELECT id, room_id, pack_id, status, current_round, current_phase, 
		       started_at, finished_at, created_at, updated_at
		FROM game_sessions
		WHERE room_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer rows.Close()

	var games []*domain.Game
	for rows.Next() {
		game := &domain.Game{
			Players: make(map[uuid.UUID]*domain.Player),
		}

		var startedAt, finishedAt sql.NullTime
		err := rows.Scan(
			&game.ID,
			&game.RoomID,
			&game.PackID,
			&game.Status,
			&game.CurrentRound,
			&game.CurrentPhase,
			&startedAt,
			&finishedAt,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}

		if startedAt.Valid {
			game.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			game.FinishedAt = &finishedAt.Time
		}

		games = append(games, game)
	}

	return games, nil
}

// GetActiveGameForUser retrieves an active game where the user is a participant
func (r *GameRepository) GetActiveGameForUser(ctx context.Context, userID uuid.UUID) (*domain.Game, error) {
	query := `
		SELECT gs.id, gs.room_id, gs.pack_id, gs.status, gs.current_round, gs.current_phase, 
		       gs.started_at, gs.finished_at, gs.created_at, gs.updated_at
		FROM game_sessions gs
		INNER JOIN game_players gp ON gs.id = gp.game_id
		WHERE gp.user_id = $1 
		  AND gp.is_active = true
		  AND gs.status NOT IN ('finished', 'cancelled', 'game_end')
		ORDER BY gs.created_at DESC
		LIMIT 1
	`

	game := &domain.Game{
		Players: make(map[uuid.UUID]*domain.Player),
	}

	var startedAt, finishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&game.ID,
		&game.RoomID,
		&game.PackID,
		&game.Status,
		&game.CurrentRound,
		&game.CurrentPhase,
		&startedAt,
		&finishedAt,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrGameNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active game for user: %w", err)
	}

	if startedAt.Valid {
		game.StartedAt = &startedAt.Time
	}
	if finishedAt.Valid {
		game.FinishedAt = &finishedAt.Time
	}

	// Load players
	players, err := r.getGamePlayers(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		game.Players[player.UserID] = player
	}

	return game, nil
}

