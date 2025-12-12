package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/domain/event"
)

func scanGame(rows *sql.Rows, game *domain.Game) error {
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
		return fmt.Errorf("failed to scan game: %w", err)
	}

	game.StartedAt = handleNullTime(startedAt)
	game.FinishedAt = handleNullTime(finishedAt)

	return nil
}

func scanGameRow(row *sql.Row, game *domain.Game) error {
	var startedAt, finishedAt sql.NullTime
	err := row.Scan(
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
		return fmt.Errorf("failed to scan game: %w", err)
	}

	game.StartedAt = handleNullTime(startedAt)
	game.FinishedAt = handleNullTime(finishedAt)

	return nil
}

func scanPlayer(rows *sql.Rows, player *domain.Player) error {
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
		return fmt.Errorf("failed to scan player: %w", err)
	}

	player.Role = domain.PlayerRole(role)
	player.LeftAt = handleNullTime(leftAt)

	return nil
}

func scanEvent(rows *sql.Rows, event *domain.GameEvent) error {
	var dataJSON []byte

	err := rows.Scan(
		&event.ID,
		&event.GameID,
		&event.EventType,
		&event.UserID,
		&event.RoundNumber,
		&event.QuestionID,
		&dataJSON,
		&event.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to scan event: %w", err)
	}

	if err := unmarshalEventData(dataJSON, event); err != nil {
		return err
	}

	return nil
}

func marshalEventData(data map[string]interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	return dataJSON, nil
}

func unmarshalEventData(dataJSON []byte, event *domain.GameEvent) error {
	if len(dataJSON) == 0 {
		return nil
	}

	event.Data = make(map[string]interface{})
	if err := json.Unmarshal(dataJSON, &event.Data); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return nil
}

func handleNullTime(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

func loadGamePlayers(ctx context.Context, db *sql.DB, gameID uuid.UUID, game *domain.Game) error {
	players, err := getGamePlayers(ctx, db, gameID)
	if err != nil {
		return err
	}

	for _, player := range players {
		game.Players[player.UserID] = player
	}

	return nil
}

func getGamePlayers(ctx context.Context, db *sql.DB, gameID uuid.UUID) ([]*domain.Player, error) {
	rows, err := db.QueryContext(ctx, querySelectGamePlayers, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game players: %w", err)
	}
	defer rows.Close()

	var players []*domain.Player
	for rows.Next() {
		player := &domain.Player{}
		if err := scanPlayer(rows, player); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

