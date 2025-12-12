package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"sigame/game/internal/domain/game"
	"sigame/game/internal/domain/pack"
	"sigame/game/internal/domain/player"
	"sigame/game/internal/domain/event"
)

func TestNewEventRepository(t *testing.T) {
	repo := NewEventRepository(nil)
	if repo == nil {
		t.Error("NewEventRepository() returned nil")
	}
	if repo.db != nil {
		t.Error("NewEventRepository() db should be nil when passed nil")
	}
}

func TestEventRepository_MethodsExist(t *testing.T) {
	repo := NewEventRepository(nil)
	ctx := context.Background()
	gameID := uuid.New()

	_ = repo.LogEvent
	_ = repo.LogEvents
	_ = repo.GetGameEvents
	_ = repo.GetEventsByType
	_ = repo.GetEventCount

	_ = ctx
	_ = gameID
}

func TestEventRepository_LogEvents_Empty(t *testing.T) {
	repo := NewEventRepository(nil)
	ctx := context.Background()

	err := repo.LogEvents(ctx, []*event.Event{})
	if err != nil {
		t.Errorf("LogEvents() with empty slice should not return error, got %v", err)
	}
}


