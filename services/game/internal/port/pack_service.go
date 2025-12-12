package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/pack"
)

type PackService interface {
	GetPackContent(ctx context.Context, packID uuid.UUID) (*pack.Pack, error)
	ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error)
}

type PackCache interface {
	GetCachedPack(ctx context.Context, packID uuid.UUID) (*pack.Pack, error)
	CachePack(ctx context.Context, p *pack.Pack) error
}

