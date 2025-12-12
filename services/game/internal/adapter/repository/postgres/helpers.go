package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sigame/game/internal/domain/player"
	"sigame/game/internal/domain/event"
	domainGame "sigame/game/internal/domain/game"
)

func scanGame(rows *sql.Rows, g *domainGame.Game) error {
	var startedAt, finishedAt sql.NullTime
	err := rows.Scan(
		&g.ID,
		&g.RoomID,
		&g.PackID,
		&g.Status,
		&g.CurrentRound,
		&g.CurrentPhase,
		&startedAt,
		&finishedAt,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to scan game: %w", err)
	}

	g.StartedAt = handleNullTime(startedAt)
	g.FinishedAt = handleNullTime(finishedAt)

	return nil
}

func scanGameRow(row *sql.Row, g *domainGame.Game) error {
	var startedAt, finishedAt sql.NullTime
	err := row.Scan(
		&g.ID,
		&g.RoomID,
		&g.PackID,
		&g.Status,
		&g.CurrentRound,
		&g.CurrentPhase,
		&startedAt,
		&finishedAt,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to scan game: %w", err)
	}

	g.StartedAt = handleNullTime(startedAt)
	g.FinishedAt = handleNullTime(finishedAt)

	return nil
}

func scanPlayer(rows *sql.Rows, p *player.Player) error {
	var role string
	var leftAt sql.NullTime

	err := rows.Scan(
		&p.UserID,
		&p.Username,
		&role,
		&p.Score,
		&p.IsActive,
		&p.JoinedAt,
		&leftAt,
	)
	if err != nil {
		return fmt.Errorf("failed to scan player: %w", err)
	}

	p.Role = player.Role(role)
	p.LeftAt = handleNullTime(leftAt)

	return nil
}

func scanEvent(rows *sql.Rows, event *event.Event) error {
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

func unmarshalEventData(dataJSON []byte, event *event.Event) error {
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

func loadGamePlayers(ctx context.Context, db *sql.DB, gameID uuid.UUID, g *domainGame.Game) error {
	players, err := getGamePlayers(ctx, db, gameID)
	if err != nil {
		return err
	}

	for _, player := range players {
		g.Players[player.UserID] = player
	}

	return nil
}

func getGamePlayers(ctx context.Context, db *sql.DB, gameID uuid.UUID) ([]*player.Player, error) {
	rows, err := db.QueryContext(ctx, querySelectGamePlayers, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get g players: %w", err)
	}
	defer rows.Close()

	var players []*player.Player
	for rows.Next() {
		player := &player.Player{}
		if err := scanPlayer(rows, player); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

