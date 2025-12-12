package port

import (
	"context"

	"github.com/sigame/game/internal/domain/event"
)

type EventLogger interface {
	LogEvent(ctx context.Context, e *event.Event) error
}

