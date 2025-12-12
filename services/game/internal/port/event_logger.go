package port

import (
	"context"

	"sigame/game/internal/domain/event"
)

type EventLogger interface {
	LogEvent(ctx context.Context, e *event.Event) error
}

