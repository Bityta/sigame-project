package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sigame/game/internal/infrastructure/logger"
	"sigame/game/proto"
)

type AuthServiceClient struct {
	conn   *grpc.ClientConn
	client proto.AuthServiceClient
}

type ValidateTokenResponse struct {
	Valid    bool
	UserID   uuid.UUID
	Username string
	Error    string
}

func NewAuthClient(address string) (*AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Auth Service: %w", err)
	}

	logger.Infof(context.Background(), "Connected to Auth Service at %s", address)

	return &AuthServiceClient{
		conn:   conn,
		client: proto.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthServiceClient) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	req := &proto.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Valid {
		return &ValidateTokenResponse{
			Valid: false,
			Error: resp.Error,
		}, nil
	}

	userID, err := uuid.Parse(resp.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in response: %w", err)
	}

	return &ValidateTokenResponse{
		Valid:    true,
		UserID:   userID,
		Username: resp.Username,
	}, nil
}

func (c *AuthServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

