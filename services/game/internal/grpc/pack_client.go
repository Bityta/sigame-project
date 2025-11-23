package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain"
	"google.golang.org/grpc"
)

// PackServiceClient defines the interface for pack service gRPC client (temporary stub until we generate proto files)
type PackServiceClient interface {
	GetPackContent(ctx context.Context, packID uuid.UUID) (*domain.Pack, error)
	ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error)
}

// PackClient is a client for Pack Service
type PackClient struct {
	address string
	conn    *grpc.ClientConn
}

// NewPackClient creates a new Pack Service client
func NewPackClient(address string, opts ...grpc.DialOption) (*PackClient, error) {
	// For now, just store the address
	// Proto generation will be done during build
	return &PackClient{
		address: address,
	}, nil
}

// Close closes the gRPC connection
func (c *PackClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetPackContent retrieves full pack content from Pack Service
// TODO: Implement after proto generation
func (c *PackClient) GetPackContent(ctx context.Context, packID uuid.UUID) (*domain.Pack, error) {
	// Temporary mock implementation
	return &domain.Pack{
		ID:     packID,
		Name:   "Mock Pack",
		Author: "System",
		Rounds: []*domain.Round{
			{
				ID:          "round1",
				RoundNumber: 1,
				Name:        "Round 1",
				Themes: []*domain.Theme{
					{
						ID:   "theme1",
						Name: "Test Theme",
						Questions: []*domain.Question{
							{
								ID:        "q1",
								Price:     100,
								Text:      "Test Question?",
								Answer:    "Test Answer",
								MediaType: "text",
								Used:      false,
							},
						},
					},
				},
			},
		},
	}, nil
}

// ValidatePackExists checks if a pack exists
// TODO: Implement after proto generation
func (c *PackClient) ValidatePackExists(ctx context.Context, packID uuid.UUID) (bool, error) {
	// For now, always return true
	return true, nil
}
