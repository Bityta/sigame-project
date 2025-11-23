package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
)

// EventRepository handles event persistence in PostgreSQL
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// LogEvent logs a game event to the database
func (r *EventRepository) LogEvent(ctx context.Context, event *domain.GameEvent) error {
	query := `
		INSERT INTO game_events (id, game_id, event_type, user_id, round_number, question_id, data, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Marshal data to JSON
	var dataJSON []byte
	var err error
	if event.Data != nil {
		dataJSON, err = json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, query,
		event.ID,
		event.GameID,
		event.EventType,
		event.UserID,
		event.RoundNumber,
		event.QuestionID,
		dataJSON,
		event.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to log event: %w", err)
	}

	return nil
}

// LogEvents logs multiple events in a batch
func (r *EventRepository) LogEvents(ctx context.Context, events []*domain.GameEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO game_events (id, game_id, event_type, user_id, round_number, question_id, data, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert all events
	for _, event := range events {
		var dataJSON []byte
		if event.Data != nil {
			dataJSON, err = json.Marshal(event.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal event data: %w", err)
			}
		}

		_, err = stmt.ExecContext(ctx,
			event.ID,
			event.GameID,
			event.EventType,
			event.UserID,
			event.RoundNumber,
			event.QuestionID,
			dataJSON,
			event.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetGameEvents retrieves all events for a game
func (r *EventRepository) GetGameEvents(ctx context.Context, gameID uuid.UUID) ([]*domain.GameEvent, error) {
	query := `
		SELECT id, game_id, event_type, user_id, round_number, question_id, data, timestamp
		FROM game_events
		WHERE game_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []*domain.GameEvent
	for rows.Next() {
		event := &domain.GameEvent{}
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
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		// Unmarshal data
		if len(dataJSON) > 0 {
			event.Data = make(map[string]interface{})
			if err := json.Unmarshal(dataJSON, &event.Data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
			}
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEventsByType retrieves events by type for a game
func (r *EventRepository) GetEventsByType(ctx context.Context, gameID uuid.UUID, eventType domain.EventType) ([]*domain.GameEvent, error) {
	query := `
		SELECT id, game_id, event_type, user_id, round_number, question_id, data, timestamp
		FROM game_events
		WHERE game_id = $1 AND event_type = $2
		ORDER BY timestamp ASC
	`

	rows, err := r.db.QueryContext(ctx, query, gameID, eventType)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var events []*domain.GameEvent
	for rows.Next() {
		event := &domain.GameEvent{}
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
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if len(dataJSON) > 0 {
			event.Data = make(map[string]interface{})
			if err := json.Unmarshal(dataJSON, &event.Data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
			}
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEventCount returns the total number of events for a game
func (r *EventRepository) GetEventCount(ctx context.Context, gameID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM game_events WHERE game_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, gameID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get event count: %w", err)
	}

	return count, nil
}

