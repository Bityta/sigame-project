package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"sigame/game/internal/domain/event"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) LogEvent(ctx context.Context, event *event.Event) error {
	dataJSON, err := marshalEventData(event.Data)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, queryInsertGameEvent,
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
		return ErrLogEvent(err)
	}

	return nil
}

func (r *EventRepository) LogEvents(ctx context.Context, events []*event.Event) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return ErrBeginTransaction(err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, queryInsertGameEvent)
	if err != nil {
		return ErrPrepareStatement(err)
	}
	defer stmt.Close()

	for _, event := range events {
		dataJSON, err := marshalEventData(event.Data)
		if err != nil {
			return err
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
			return ErrInsertEvent(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return ErrCommitTransaction(err)
	}

	return nil
}

func (r *EventRepository) GetGameEvents(ctx context.Context, gameID uuid.UUID) ([]*event.Event, error) {
	rows, err := r.db.QueryContext(ctx, querySelectGameEvents, gameID)
	if err != nil {
		return nil, ErrGetEvents(err)
	}
	defer rows.Close()

	var events []*event.Event
	for rows.Next() {
		event := &event.Event{}
		if err := scanEvent(rows, event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetEventsByType(ctx context.Context, gameID uuid.UUID, eventType event.Type) ([]*event.Event, error) {
	rows, err := r.db.QueryContext(ctx, querySelectEventsByType, gameID, eventType)
	if err != nil {
		return nil, ErrGetEvents(err)
	}
	defer rows.Close()

	var events []*event.Event
	for rows.Next() {
		event := &event.Event{}
		if err := scanEvent(rows, event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetEventCount(ctx context.Context, gameID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryCountGameEvents, gameID).Scan(&count)
	if err != nil {
		return 0, ErrGetEventCount(err)
	}

	return count, nil
}

